package main

import (
	"log"
	"net/http"
	"os"

	"inspector-gadget-management/backend/internal/gadget"
	"inspector-gadget-management/backend/internal/handler"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize gadget client
	gadgetClient := gadget.NewClient()

	// Initialize handler
	h := handler.NewHandler(gadgetClient)

	// Setup router
	r := mux.NewRouter()
	h.RegisterRoutes(r)

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
