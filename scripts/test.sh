#!/bin/bash

# Test script for language detection service
# This script runs all unit tests with coverage

set -e

echo "🧪 Running Language Detection Service Tests"
echo "=============================================="

# Change to project root
cd "$(dirname "$0")/.."

echo "📦 Running go mod tidy..."
go mod tidy

echo "🔍 Running all tests..."
go test -v ./...

echo ""
echo "📊 Test Coverage Report:"
echo "========================"
go test -cover ./...

echo ""
echo "✅ All tests completed successfully!"
echo ""
echo "💡 To run tests with detailed coverage:"
echo "   go test -coverprofile=coverage.out ./..."
echo "   go tool cover -html=coverage.out"
