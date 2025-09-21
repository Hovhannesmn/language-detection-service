# Language Detection gRPC Service

A high-performance gRPC service for language detection using AWS Comprehend AI with intelligent fallback pattern matching.

## Features

- **gRPC Protocol**: High-performance language detection service
- **AWS Comprehend AI**: Accurate language detection using AWS AI
- **Intelligent Fallback**: Pattern-based detection when AWS is unavailable
- **Batch Processing**: Support for batch language detection requests
- **Graceful Shutdown**: Proper signal handling for clean shutdowns
- **Multiple Languages**: Supports 13+ languages (English, Spanish, French, German, Italian, Portuguese, Russian, Japanese, Korean, Chinese, Arabic, Hindi, and more)

## gRPC API

```protobuf
rpc DetectLanguage(DetectLanguageRequest) returns (DetectLanguageResponse);
```

## Configuration

- **Server Address**: `0.0.0.0:6011`
- **AWS Region**: `us-east-1`
- **Max Text Length**: `5000` characters
- **Min Confidence**: `0.10` (10%)

## ⚠️ IMPORTANT: AWS Configuration Required

# **TO USE AWS COMPREHEND AI, SET YOUR CREDENTIALS:**

```bash
export AWS_ACCESS_KEY_ID=your_access_key_here
export AWS_SECRET_ACCESS_KEY=your_secret_access_key_here
export AWS_REGION=us-east-1
```

**Required AWS Permissions:**
- `comprehend:DetectDominantLanguage`
- `comprehend:BatchDetectDominantLanguage`

> **Note**: Without AWS credentials, the service will automatically fall back to pattern-based detection.

## Prerequisites

- Go 1.24.2 or higher
- AWS Account with Comprehend access (optional - fallback available)

## Installation & Running

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Run the Service
```bash
go run cmd/server/main.go
```

### 3. Using Docker
```bash
docker build -t language-detection-service .
docker run -p 6011:6011 language-detection-service
```

## Fallback Behavior

The service automatically falls back to pattern-based detection when:
- AWS credentials not configured
- AWS Comprehend unavailable
- Network connectivity issues

## Service Details

- **Address**: `0.0.0.0:6011`
- **Protocol**: gRPC
- **Service Name**: `language_detection.LanguageDetectionService`

## Testing

```bash
# Run all tests
go test ./...

# Test with gRPC client
grpcurl -plaintext -d '{"text": "Hello, world!"}' localhost:6011 language_detection.LanguageDetectionService/DetectLanguage
```

## Example Usage

### Using grpcurl
```bash
grpcurl -plaintext -d '{"text": "Bonjour le monde!"}' \
  localhost:6011 \
  language_detection.LanguageDetectionService/DetectLanguage
```

### Using Go Client
```go
import pb "github.com/Hovhannesmn/ld_proto/pb"

conn, err := grpc.Dial("localhost:6011", grpc.WithInsecure())
client := pb.NewLanguageDetectionServiceClient(conn)

response, err := client.DetectLanguage(context.Background(), &pb.DetectLanguageRequest{
    Text: "Hello, world!",
})
```

## Graceful Shutdown

The service supports graceful shutdown via SIGINT (Ctrl+C) or SIGTERM signals with proper cleanup of existing connections.

## License

MIT License