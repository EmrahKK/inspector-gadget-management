package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"inspector-gadget-management/backend/internal/gadget"
	"inspector-gadget-management/backend/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Handler manages HTTP and WebSocket handlers
type Handler struct {
	gadgetClient *gadget.Client
	upgrader     websocket.Upgrader
	wsClients    map[string]*WSClient
	mu           sync.RWMutex
}

// WSClient represents a WebSocket client
type WSClient struct {
	SessionID string
	Conn      *websocket.Conn
	Send      chan []byte
}

// NewHandler creates a new handler
func NewHandler(gadgetClient *gadget.Client) *Handler {
	return &Handler{
		gadgetClient: gadgetClient,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		wsClients: make(map[string]*WSClient),
	}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	// API routes
	r.HandleFunc("/api/gadgets", h.ListGadgets).Methods("GET")
	r.HandleFunc("/api/sessions", h.ListSessions).Methods("GET")
	r.HandleFunc("/api/sessions", h.StartSession).Methods("POST")
	r.HandleFunc("/api/sessions/{sessionId}", h.StopSession).Methods("DELETE")

	// WebSocket route
	r.HandleFunc("/ws/{sessionId}", h.HandleWebSocket)
}

// ListGadgets returns available gadgets
func (h *Handler) ListGadgets(w http.ResponseWriter, r *http.Request) {
	gadgets := []map[string]interface{}{
		{
			"type":        models.GadgetTraceTCP,
			"name":        "Trace TCP",
			"description": "Trace TCP connections",
			"category":    "trace",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gadgets)
}

// ListSessions returns all active sessions
func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	sessions := h.gadgetClient.ListSessions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// StartSession starts a new gadget session
func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	var req models.GadgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	sessionID := uuid.New().String()

	// Use background context so gadget continues running after HTTP request completes
	session, err := h.gadgetClient.RunGadget(context.Background(), req, sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start gadget: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.GadgetSession{
		ID:        session.ID,
		Type:      session.Type,
		Namespace: session.Namespace,
		PodName:   session.PodName,
		Status:    session.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// StopSession stops a running gadget session
func (h *Handler) StopSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	if err := h.gadgetClient.StopGadget(sessionID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to stop session: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleWebSocket handles WebSocket connections for real-time gadget output
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	session, exists := h.gadgetClient.GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &WSClient{
		SessionID: sessionID,
		Conn:      conn,
		Send:      make(chan []byte, 256),
	}

	h.mu.Lock()
	h.wsClients[sessionID] = client
	h.mu.Unlock()

	// Start goroutines for reading and writing
	go h.wsWriter(client)
	go h.wsReader(client)
	go h.forwardGadgetOutput(session, client)
}

// wsWriter writes messages to WebSocket
func (h *Handler) wsWriter(client *WSClient) {
	defer func() {
		client.Conn.Close()
		h.mu.Lock()
		delete(h.wsClients, client.SessionID)
		h.mu.Unlock()
	}()

	for {
		message, ok := <-client.Send
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}

// wsReader reads messages from WebSocket (for keepalive)
func (h *Handler) wsReader(client *WSClient) {
	defer client.Conn.Close()

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
	}
}

// forwardGadgetOutput forwards gadget output to WebSocket client
func (h *Handler) forwardGadgetOutput(session *gadget.Session, client *WSClient) {
	for {
		select {
		case output, ok := <-session.OutputCh:
			if !ok {
				// Channel closed, session ended
				message := map[string]interface{}{
					"type":   "session_ended",
					"status": session.Status,
				}
				if data, err := json.Marshal(message); err == nil {
					client.Send <- data
				}
				close(client.Send)
				return
			}

			// Forward output to WebSocket
			if data, err := json.Marshal(output); err == nil {
				select {
				case client.Send <- data:
				default:
					// Client send buffer full, skip message
				}
			}

		case err, ok := <-session.ErrorCh:
			if !ok {
				continue
			}

			// Forward error to WebSocket
			errorMsg := map[string]interface{}{
				"type":    "error",
				"message": err.Error(),
			}
			if data, err := json.Marshal(errorMsg); err == nil {
				select {
				case client.Send <- data:
				default:
					// Client send buffer full, skip message
				}
			}
		}
	}
}
