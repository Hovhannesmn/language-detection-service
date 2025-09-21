# Multi-stage Dockerfile for Language Detection Service
# Stage 1: Build stage
FROM golang:1.24.2-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Copy ld-proto directory (needed for replace directive)
COPY ld-proto ./ld-proto

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o language-detection-service cmd/server/main.go

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/language-detection-service .

# Copy ld-proto directory (if needed for runtime)
COPY --from=builder /app/ld-proto ./ld-proto

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 6011


# Set environment variables with defaults
ENV SERVER_ADDRESS=0.0.0.0
ENV SERVER_PORT=6011
ENV AWS_REGION=us-east-1
ENV USE_AWS_COMPREHEND=true
ENV MAX_TEXT_LENGTH=5000
ENV MIN_CONFIDENCE_THRESHOLD=0.1
ENV SERVICE_VERSION=1.0.0
ENV MODEL_VERSION=1.0.0
ENV SHUTDOWN_TIMEOUT_SECONDS=30

# Run the application
CMD ["./language-detection-service"]
