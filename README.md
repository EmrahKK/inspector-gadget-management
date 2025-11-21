# Inspector Gadget Management

A web-based management interface for [Inspektor Gadget](https://github.com/inspektor-gadget/inspektor-gadget), providing real-time monitoring and tracing of Kubernetes workloads using eBPF.

## Architecture

```
┌────────────────────┐
│   Web Frontend      │ (React + TypeScript)
│   - Visualizations │
│   - Run Gadgets     │
└────────┬───────────┘
         │ WebSocket/REST
┌────────▼────────┐
│  Backend API    │ (Go)
│  - Gadget mgmt  │
│  - Output parse │
└────────┬────────┘
         │ kubectl-gadget CLI
┌────────▼────────┐
│ Inspektor-Gadget│ (DaemonSet)
│  on K8s cluster │
└─────────────────┘
```

## Features

- **Web-based UI**: Modern React interface for managing gadgets
- **Real-time Output**: WebSocket-based streaming of gadget events
- **Multiple Gadgets**: Support for popular gadgets:
  - `trace_exec`: Trace process execution
  - `trace_tcp`: Trace TCP connections
  - `snapshot_process`: Snapshot running processes
- **Namespace Filtering**: Filter gadgets by namespace or pod
- **Session Management**: Start, stop, and monitor multiple gadget sessions

## Prerequisites

- Kubernetes cluster (k3s, k3d, minikube, etc.)
- [Inspektor Gadget](https://www.inspektor-gadget.io/docs/latest/quick-start/) installed on the cluster
- Podman or Docker for building images
- kubectl configured to access your cluster

## Quick Start

### 1. Install Inspektor Gadget

If you haven't already installed Inspektor Gadget on your cluster:

```bash
kubectl gadget deploy
```

Verify the installation:

```bash
kubectl gadget version
```

### 2. Build Container Images

The build script automatically detects whether you're using Podman or Docker:

```bash
./build.sh
```

For k3s, import the images:

```bash
# Using Podman
podman save gadget-backend:latest | sudo k3s ctr images import -
podman save gadget-frontend:latest | sudo k3s ctr images import -

# Or using the Makefile
make import-k3s
```

For k3d:

```bash
k3d image import gadget-backend:latest gadget-frontend:latest -c mycluster
```

For minikube:

```bash
minikube image load gadget-backend:latest
minikube image load gadget-frontend:latest
```

### 3. Deploy to Kubernetes

```bash
./deploy.sh
```

### 4. Access the Application

**Option 1: Port Forward**

```bash
kubectl port-forward -n gadget-management svc/frontend 3000:80
```

Then open http://localhost:3000

**Option 2: NodePort**

Access via NodePort (default: 30080):

```bash
# For k3s/k3d
http://localhost:30080

# For other clusters, get the node IP
kubectl get nodes -o wide
# Then access: http://<NODE_IP>:30080
```

**Option 3: Ingress (Optional)**

Deploy the Ingress resource:

```bash
kubectl apply -f k8s/ingress.yaml
```

Add to `/etc/hosts`:

```
127.0.0.1 gadget.local
```

Access: http://gadget.local

## Usage

### Starting a Gadget

1. Select a gadget type from the dropdown
2. Optionally specify a namespace and/or pod name
3. Click "Start Gadget"
4. View real-time events in the output panel

### Managing Sessions

- **View Active Sessions**: See all running gadget sessions in the sidebar
- **Switch Sessions**: Click on a session to view its output
- **Stop Session**: Click the "Stop" button on any session

### Example: Trace Process Execution

1. Start a `trace_exec` gadget
2. In another terminal, create some processes in your cluster:
   ```bash
   kubectl run test-pod --image=busybox -- sh -c "while true; do echo hello; sleep 1; done"
   ```
3. See the exec events appear in real-time in the UI

## Project Structure

```
.
├── backend/                    # Go backend service
│   ├── cmd/server/            # Main application entry point
│   ├── internal/
│   │   ├── gadget/           # Gadget client implementation
│   │   ├── handler/          # HTTP and WebSocket handlers
│   │   └── models/           # Data models
│   ├── Dockerfile
│   └── go.mod
├── frontend/                  # React frontend
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── services/         # API services
│   │   └── types.ts          # TypeScript types
│   ├── Dockerfile
│   ├── nginx.conf            # Nginx configuration
│   └── package.json
├── k8s/                       # Kubernetes manifests
│   ├── namespace.yaml
│   ├── backend-rbac.yaml
│   ├── backend-deployment.yaml
│   ├── frontend-deployment.yaml
│   └── ingress.yaml
├── build.sh                   # Build Docker images
├── deploy.sh                  # Deploy to Kubernetes
└── README.md
```

## API Endpoints

### REST API

- `GET /api/gadgets` - List available gadgets
- `GET /api/sessions` - List active sessions
- `POST /api/sessions` - Start a new gadget session
- `DELETE /api/sessions/{sessionId}` - Stop a session
- `GET /health` - Health check

### WebSocket

- `WS /ws/{sessionId}` - Stream gadget output for a session

## Container Runtime Notes

This project supports both **Podman** and **Docker**. The build script automatically detects which one is available on your system.

**Platform Support:** Images are built for **linux/amd64** by default, which is compatible with most Kubernetes clusters.

### Using Podman

Podman is a daemonless container engine that can run containers without root privileges:

```bash
# Build images for linux/amd64
./build.sh  # Automatically uses podman if available and builds for linux/amd64

# Import to k3s
make import-k3s

# Or manually
podman save gadget-backend:latest | sudo k3s ctr images import -
podman save gadget-frontend:latest | sudo k3s ctr images import -
```

**Note for Apple Silicon (M1/M2/M3):** The build script automatically uses `--platform linux/amd64` to ensure compatibility with AMD64-based Kubernetes nodes.

### Using Docker

If you prefer Docker, the scripts will automatically use it if Podman is not available:

```bash
# Build images for linux/amd64
./build.sh  # Automatically uses docker if podman not found

# Import to k3s
docker save gadget-backend:latest | sudo k3s ctr images import -
docker save gadget-frontend:latest | sudo k3s ctr images import -
```

## Development

### Backend Development

```bash
cd backend
go mod download
go run cmd/server/main.go
```

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at http://localhost:3000 with hot reload.

## Configuration

### Backend Environment Variables

- `PORT`: Server port (default: 8080)

### Frontend Environment Variables

Create a `.env` file in the `frontend` directory:

- `VITE_API_URL`: Backend API URL (default: /api)
- `VITE_WS_URL`: WebSocket URL (default: window.location.host)

## Troubleshooting

### Pods not starting

Check the logs:

```bash
kubectl logs -n gadget-management -l app=gadget-backend
kubectl logs -n gadget-management -l app=gadget-frontend
```

### kubectl-gadget not working

Ensure Inspektor Gadget is properly installed:

```bash
kubectl get pods -n gadget
kubectl gadget version
```

### Permission errors

The backend needs permissions to run kubectl-gadget. Check RBAC:

```bash
kubectl get clusterrole gadget-backend-role
kubectl get clusterrolebinding gadget-backend-binding
```

### WebSocket connection failed

1. Check that the backend is running
2. Verify network policies allow WebSocket connections
3. Check browser console for errors

## Clean Up

Remove the deployment:

```bash
kubectl delete namespace gadget-management
# or using Makefile
make clean
```

Remove container images:

```bash
# Using Podman
podman rmi gadget-backend:latest gadget-frontend:latest

# Using Docker
docker rmi gadget-backend:latest gadget-frontend:latest
```

## Future Enhancements

- [ ] Support for more gadget types
- [ ] Export gadget output (JSON, CSV)
- [ ] Authentication and authorization
- [ ] Multi-cluster support
- [ ] Persistent storage for historical data
- [ ] Advanced filtering and search
- [ ] Custom gadget parameters
- [ ] Real-time statistics and charts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Resources

- [Inspektor Gadget Documentation](https://www.inspektor-gadget.io/docs/)
- [eBPF Introduction](https://ebpf.io/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
