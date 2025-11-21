import React, { useEffect, useState } from 'react';
import { GadgetSelector } from './components/GadgetSelector';
import { SessionList } from './components/SessionList';
import { GadgetOutput } from './components/GadgetOutput';
import { GadgetRequest, GadgetSession, GadgetOutput as GadgetOutputType } from './types';
import { api } from './services/api';

function App() {
  const [sessions, setSessions] = useState<GadgetSession[]>([]);
  const [activeSessionId, setActiveSessionId] = useState<string | undefined>();
  const [outputs, setOutputs] = useState<GadgetOutputType[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [sidebarCollapsed, setSidebarCollapsed] = useState<boolean>(false);

  useEffect(() => {
    loadSessions();
    const interval = setInterval(loadSessions, 5000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (activeSessionId) {
      connectWebSocket(activeSessionId);
    }
    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [activeSessionId]);

  const loadSessions = async () => {
    try {
      const data = await api.getSessions();
      setSessions(data);
    } catch (error) {
      console.error('Failed to load sessions:', error);
    }
  };

  const connectWebSocket = (sessionId: string) => {
    if (ws) {
      ws.close();
    }

    setOutputs([]);
    const wsUrl = api.getWebSocketUrl(sessionId);
    const websocket = new WebSocket(wsUrl);

    websocket.onopen = () => {
      console.log('WebSocket connected');
      setError(null);
    };

    websocket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);

        if (message.type === 'error') {
          setError(message.message);
          return;
        }

        if (message.type === 'session_ended') {
          console.log('Session ended:', message.status);
          loadSessions();
          return;
        }

        // Regular gadget output
        if (message.data) {
          setOutputs((prev) => [...prev, message]);
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
      setError('WebSocket connection error');
    };

    websocket.onclose = () => {
      console.log('WebSocket disconnected');
    };

    setWs(websocket);
  };

  const handleStartGadget = async (request: GadgetRequest) => {
    try {
      setError(null);
      const session = await api.startSession(request);
      setSessions((prev) => [...prev, session]);
      setActiveSessionId(session.id);
    } catch (error: any) {
      setError(error.response?.data || 'Failed to start gadget');
      console.error('Failed to start gadget:', error);
    }
  };

  const handleStopSession = async (sessionId: string) => {
    try {
      await api.stopSession(sessionId);
      setSessions((prev) => prev.filter((s) => s.id !== sessionId));
      if (activeSessionId === sessionId) {
        setActiveSessionId(undefined);
        setOutputs([]);
      }
    } catch (error) {
      console.error('Failed to stop session:', error);
    }
  };

  const handleSelectSession = (sessionId: string) => {
    setActiveSessionId(sessionId);
  };

  return (
    <div style={styles.app}>
      <header style={styles.header}>
        <h1 style={styles.headerTitle}>Inspector Gadget Management</h1>
        <p style={styles.headerSubtitle}>
          Monitor and trace your Kubernetes workloads with eBPF
        </p>
      </header>

      {error && (
        <div style={styles.errorBanner}>
          <strong>Error:</strong> {error}
          <button style={styles.errorClose} onClick={() => setError(null)}>
            ×
          </button>
        </div>
      )}

      <div style={styles.container}>
        <div style={{
          ...styles.sidebar,
          ...(sidebarCollapsed ? styles.sidebarCollapsed : {}),
        }}>
          {!sidebarCollapsed && (
            <>
              <GadgetSelector
                onStartGadget={handleStartGadget}
                disabled={false}
              />
              <SessionList
                sessions={sessions}
                activeSessionId={activeSessionId}
                onSelectSession={handleSelectSession}
                onStopSession={handleStopSession}
              />
            </>
          )}
          <button
            style={{
              ...styles.toggleButton,
              ...(sidebarCollapsed ? styles.toggleButtonCollapsed : {}),
            }}
            onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
            title={sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
          >
            {sidebarCollapsed ? '▶' : '◀'}
          </button>
        </div>

        <div style={styles.main}>
          <GadgetOutput outputs={outputs} sessionId={activeSessionId} />
        </div>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  app: {
    minHeight: '100vh',
    backgroundColor: '#e9ecef',
  },
  header: {
    backgroundColor: '#007bff',
    color: 'white',
    padding: '20px 40px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
  },
  headerTitle: {
    margin: 0,
    fontSize: '28px',
    fontWeight: 'bold',
  },
  headerSubtitle: {
    margin: '5px 0 0 0',
    fontSize: '14px',
    opacity: 0.9,
  },
  errorBanner: {
    backgroundColor: '#f8d7da',
    color: '#721c24',
    padding: '12px 20px',
    margin: '20px 40px',
    borderRadius: '4px',
    border: '1px solid #f5c6cb',
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  errorClose: {
    backgroundColor: 'transparent',
    border: 'none',
    fontSize: '24px',
    cursor: 'pointer',
    color: '#721c24',
    padding: '0 5px',
  },
  container: {
    display: 'flex',
    gap: '20px',
    padding: '20px 40px',
    maxWidth: '1600px',
    margin: '0 auto',
  },
  sidebar: {
    flex: '0 0 400px',
    position: 'relative',
    transition: 'all 0.3s ease',
  },
  sidebarCollapsed: {
    flex: '0 0 60px',
  },
  toggleButton: {
    position: 'absolute',
    top: '10px',
    right: '10px',
    backgroundColor: '#007bff',
    color: 'white',
    border: 'none',
    borderRadius: '50%',
    width: '36px',
    height: '36px',
    cursor: 'pointer',
    fontSize: '14px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
    transition: 'all 0.2s ease',
    zIndex: 10,
  },
  toggleButtonCollapsed: {
    right: '12px',
  },
  main: {
    flex: 1,
    minWidth: 0,
  },
};

export default App;
