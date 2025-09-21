#!/bin/bash

# Generate protobuf code for the language detection service

set -e

echo "Generating protobuf code..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Please install Protocol Buffers compiler:"
    echo "  macOS: brew install protobuf"
    echo "  Ubuntu: apt-get install protobuf-compiler"
    echo "  Or download from: https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directory if it doesn't exist
mkdir -p proto

# Generate Go code from protobuf
echo "Generating Go code from proto/language_detection.proto..."
mkdir -p pb-service
protoc --go_out=pb-service --go_opt=paths=source_relative \
    --go-grpc_out=pb-service --go-grpc_opt=paths=source_relative \
    proto/language_detection.proto

echo "Protobuf code generation completed!"
echo "Generated files:"
echo "  - pb-service/proto/language_detection.pb.go"
echo "  - pb-service/proto/language_detection_grpc.pb.go"
