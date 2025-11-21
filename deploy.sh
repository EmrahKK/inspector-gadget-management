#!/bin/bash
set -e

echo "Deploying Inspector Gadget Management to Kubernetes..."

# Apply Kubernetes manifests
echo "Creating namespace..."
kubectl apply -f k8s/namespace.yaml

echo "Creating RBAC resources..."
kubectl apply -f k8s/backend-rbac.yaml

echo "Deploying backend..."
kubectl apply -f k8s/backend-deployment.yaml

echo "Deploying frontend..."
kubectl apply -f k8s/frontend-deployment.yaml

echo "Deploying ingress..."
kubectl apply -f k8s/ingress.yaml

echo ""
echo "Deployment completed!"
echo ""
echo "Checking deployment status..."
kubectl get pods -n gadget-management

echo ""
echo "Services:"
kubectl get svc -n gadget-management

echo ""
echo "Access the application:"
echo "  NodePort: http://<node-ip>:30080"
echo "  Port-forward: kubectl port-forward -n gadget-management svc/frontend 3000:80"
echo "               Then open http://localhost:3000"
echo ""
echo "To deploy Ingress (optional):"
echo "  kubectl apply -f k8s/ingress.yaml"
