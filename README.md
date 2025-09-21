# Language Detection gRPC Service

A standalone gRPC service for language detection using AWS Comprehend with fallback pattern matching.

## Features

- **gRPC API**: High-performance language detection service
- **AWS Comprehend Integration**: Uses AWS Comprehend for accurate language detection
- **Fallback Support**: Pattern-based detection when AWS is unavailable
- **Batch Processing**: Support for batch language detection
- **Health Checks**: Built-in health monitoring
- **Graceful Shutdown**: Proper signal handling for clean shutdowns

## API Endpoints

### Single Language Detection
```protobuf
rpc DetectLanguage(DetectLanguageRequest) returns (DetectLanguageResponse);
```

### Batch Language Detection
```protobuf
rpc DetectLanguageBatch(DetectLanguageBatchRequest) returns (DetectLanguageBatchResponse);
```

### Health Check
```protobuf
rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
```

## Configuration

The service supports the following configuration options:

- **Server Address**: `0.0.0.0:8090` (default)
- **AWS Comprehend**: Enabled by default
- **AWS Region**: `us-east-1` (default)
- **Max Text Length**: `5000` characters
- **Min Confidence**: `0.10` (10%)
- **Supported Languages**: Multiple languages including English, Spanish, French, German, etc.

## Running the Service

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Generate protobuf files**:
   ```bash
   chmod +x scripts/generate_proto.sh
   ./scripts/generate_proto.sh
   ```

3. **Run the service**:
   ```bash
   go run cmd/server/main.go
   ```

4. **Optional: Set AWS credentials** (if using AWS Comprehend):
   ```bash
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   export AWS_REGION=us-east-1
   ```

## Service Details

- **Address**: `localhost:8090`
- **Protocol**: gRPC
- **Service Name**: `language_detection.LanguageDetectionService`
- **Health Service**: `grpc.health.v1.Health`

## Graceful Shutdown

The service supports graceful shutdown via SIGINT (Ctrl+C) or SIGTERM signals. When a shutdown signal is received:

1. Context is cancelled
2. gRPC server stops accepting new connections
3. Existing requests are completed
4. Server shuts down cleanly

## Example Usage

The service can be used with any gRPC client. See the `language-detection-client` project for a complete client example.
