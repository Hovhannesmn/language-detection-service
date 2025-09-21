# Language Detection Service Makefile

# Build the service
build:
	go build -o language-detection-service cmd/server/main.go

# Run the service
run:
	go run cmd/server/main.go

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run tests with detailed coverage report
test-coverage-detail:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run tests using the test script
test-script:
	./scripts/test.sh

# Clean build artifacts
clean:
	rm -f language-detection-service
	rm -f coverage.out
	rm -rf coverage.html

# Clean generated proto files
clean-proto:
	rm -rf pb-service/proto/*.pb.go

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Docker commands
docker-build:
	docker build -t language-detection-service .

docker-run:
	docker run -p 6011:6011 language-detection-service

docker-compose-up:
	docker-compose up --build

docker-compose-up-detached:
	docker-compose up -d --build

docker-compose-down:
	docker-compose down

docker-compose-logs:
	docker-compose logs -f

# AWS mode (set USE_AWS_COMPREHEND=true)
docker-compose-aws:
	USE_AWS_COMPREHEND=true docker-compose up --build

docker-compose-aws-detached:
	USE_AWS_COMPREHEND=true docker-compose up -d --build

# Run all checks (format, lint, test)
check: fmt lint test

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Build & Run:"
	@echo "  build              - Build the service binary"
	@echo "  run                - Run the service"
	@echo ""
	@echo "Testing:"
	@echo "  test               - Run all tests"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-coverage-detail - Generate detailed coverage report"
	@echo "  test-script        - Run tests using the test script"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo ""
	@echo "Docker Compose:"
	@echo "  docker-compose-up  - Start services (fallback mode)"
	@echo "  docker-compose-up-detached - Start services in background"
	@echo "  docker-compose-down - Stop services"
	@echo "  docker-compose-logs - View service logs"
	@echo "  docker-compose-aws - Start with AWS Comprehend enabled"
	@echo "  docker-compose-aws-detached - Start AWS mode in background"
	@echo ""
	@echo "Development:"
	@echo "  clean              - Clean build artifacts"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter"
	@echo "  proto              - Generate protobuf files"
	@echo "  clean-proto        - Clean generated proto files"
	@echo "  deps               - Install dependencies"
	@echo "  check              - Run all checks (format, lint, test)"
	@echo ""
	@echo "  help               - Show this help message"