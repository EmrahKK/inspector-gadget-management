#!/bin/bash
set -e

# Check if version parameter is provided
if [ -z "$1" ]; then
    echo "Error: Version parameter is required"
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION="$1"
PLATFORM="linux/amd64"
TARGETARCH="amd64"

# Detect container runtime (prefer docker for buildx support)
if command -v docker &> /dev/null; then
    CONTAINER_RUNTIME="docker"
elif command -v podman &> /dev/null; then
    CONTAINER_RUNTIME="podman"
else
    echo "Error: Neither docker nor podman found. Please install one of them."
    exit 1
fi

echo "Using container runtime: $CONTAINER_RUNTIME"
echo "Building for platform: $PLATFORM"
echo "Version: $VERSION"
echo ""

# Docker Hub username
DOCKER_USERNAME="emrahkk"

# Image names with version
BACKEND_IMAGE="$DOCKER_USERNAME/gadget-backend:$VERSION"
FRONTEND_IMAGE="$DOCKER_USERNAME/gadget-frontend:$VERSION"

# Ensure buildx builder exists and is using docker-container driver
if [ "$CONTAINER_RUNTIME" = "docker" ]; then
    echo "Setting up buildx builder..."
    # Create builder if it doesn't exist
    if ! docker buildx inspect multiplatform &> /dev/null; then
        docker buildx create --name multiplatform --driver docker-container --use
    else
        docker buildx use multiplatform
    fi
fi

# Build and push backend
echo "Building and pushing backend image for $PLATFORM..."
if [ "$CONTAINER_RUNTIME" = "docker" ]; then
    docker buildx build --platform $PLATFORM \
        --build-arg TARGETARCH=$TARGETARCH \
        -t $BACKEND_IMAGE \
        --push \
        ./backend
else
    $CONTAINER_RUNTIME build --platform $PLATFORM \
        --build-arg TARGETARCH=$TARGETARCH \
        -t gadget-backend:$VERSION \
        ./backend
    $CONTAINER_RUNTIME tag gadget-backend:$VERSION $BACKEND_IMAGE
    echo "Pushing backend..."
    $CONTAINER_RUNTIME push $BACKEND_IMAGE
fi

# Build and push frontend
echo "Building and pushing frontend image for $PLATFORM..."
if [ "$CONTAINER_RUNTIME" = "docker" ]; then
    docker buildx build --platform $PLATFORM \
        -t $FRONTEND_IMAGE \
        --push \
        ./frontend
else
    $CONTAINER_RUNTIME build --platform $PLATFORM \
        -t gadget-frontend:$VERSION \
        ./frontend
    $CONTAINER_RUNTIME tag gadget-frontend:$VERSION $FRONTEND_IMAGE
    echo "Pushing frontend..."
    $CONTAINER_RUNTIME push $FRONTEND_IMAGE
fi

echo ""
echo "Images pushed successfully!"
echo "  - $BACKEND_IMAGE"
echo "  - $FRONTEND_IMAGE"
echo ""
echo "To use these images in Kubernetes:"
echo "  Update your deployment files to use:"
echo "    Backend: $BACKEND_IMAGE"
echo "    Frontend: $FRONTEND_IMAGE"
