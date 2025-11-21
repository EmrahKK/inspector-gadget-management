#!/bin/bash
set -e

# Check if version parameter is provided
if [ -z "$1" ]; then
    echo "Error: Version parameter is required"
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION=$1

# Check if docker is available (required for buildx)
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is required for cross-platform builds"
    echo "Please install Docker Desktop for Mac"
    exit 1
fi

echo "Using Docker with buildx for cross-platform build"
echo "Building for linux/amd64 platform with version: $VERSION"

# Set platform and architecture
PLATFORM="linux/amd64"
TARGETARCH="amd64"

# Docker Hub username
DOCKER_USERNAME="emrahkk"

# Image names with version tag
BACKEND_IMAGE="$DOCKER_USERNAME/gadget-backend:$VERSION"
FRONTEND_IMAGE="$DOCKER_USERNAME/gadget-frontend:$VERSION"

# Ensure buildx is available
echo ""
echo "Checking docker buildx..."
docker buildx version

# Create builder instance if it doesn't exist
if ! docker buildx ls | grep -q "multiplatform"; then
    echo "Creating multiplatform builder..."
    docker buildx create --name multiplatform --use
else
    echo "Using existing multiplatform builder..."
    docker buildx use multiplatform
fi

# Build and push backend
echo ""
echo "Building and pushing backend image for $PLATFORM..."
docker buildx build --platform $PLATFORM \
    --build-arg TARGETARCH=$TARGETARCH \
    -t $BACKEND_IMAGE \
    --push \
    ./backend

# Build and push frontend
echo ""
echo "Building and pushing frontend image for $PLATFORM..."
docker buildx build --platform $PLATFORM \
    -t $FRONTEND_IMAGE \
    --push \
    ./frontend

echo ""
echo "======================================"
echo "Successfully built and pushed images:"
echo "  - $BACKEND_IMAGE"
echo "  - $FRONTEND_IMAGE"
echo "======================================"
