#!/bin/bash

# Docker build script for Language Detection Service

set -e

echo "🐳 Building Language Detection Service Docker Image"
echo "=================================================="

# Change to project root
cd "$(dirname "$0")/.."

# Default values
IMAGE_NAME="language-detection-service"
TAG="latest"
PLATFORM="linux/amd64"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -t|--tag)
      TAG="$2"
      shift 2
      ;;
    -p|--platform)
      PLATFORM="$2"
      shift 2
      ;;
    -n|--name)
      IMAGE_NAME="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [OPTIONS]"
      echo "Options:"
      echo "  -t, --tag TAG        Docker image tag (default: latest)"
      echo "  -p, --platform       Target platform (default: linux/amd64)"
      echo "  -n, --name NAME      Docker image name (default: language-detection-service)"
      echo "  -h, --help           Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

FULL_IMAGE_NAME="${IMAGE_NAME}:${TAG}"

echo "📦 Building image: ${FULL_IMAGE_NAME}"
echo "🎯 Platform: ${PLATFORM}"

# Build the Docker image
docker build \
  --platform "${PLATFORM}" \
  --tag "${FULL_IMAGE_NAME}" \
  --file Dockerfile \
  .

echo ""
echo "✅ Docker image built successfully!"
echo "📋 Image details:"
docker images "${FULL_IMAGE_NAME}"

echo ""
echo "🚀 To run the container:"
echo "   docker run -p 6011:6011 ${FULL_IMAGE_NAME}"

echo ""
echo "🏥 To test gRPC service:"
echo "   Use a gRPC client to connect to localhost:6011"

echo ""
echo "📊 To view container logs:"
echo "   docker logs <container_id>"
