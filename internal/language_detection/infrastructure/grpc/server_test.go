package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/Hovhannesmn/ld_proto/pb"
	"language-detection-service/internal/language_detection/domain"
)

// MockLanguageDetectionService is a mock implementation of LanguageDetectionService
type MockLanguageDetectionService struct {
	response *domain.LanguageDetectionResponse
	err      error
}

func (m *MockLanguageDetectionService) DetectLanguage(ctx context.Context, request *domain.LanguageDetectionRequest) (*domain.LanguageDetectionResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func TestNewServer(t *testing.T) {
	mockService := &MockLanguageDetectionService{}

	server := NewServer(mockService)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.service != mockService {
		t.Errorf("Expected service to be %v, got %v", mockService, server.service)
	}

	if server.server == nil {
		t.Error("Expected gRPC server to be initialized, got nil")
	}

	if server.healthServer == nil {
		t.Error("Expected health server to be initialized, got nil")
	}

	if server.shutdownTimeout != 30*time.Second {
		t.Errorf("Expected shutdown timeout 30s, got %v", server.shutdownTimeout)
	}
}

func TestNewServer_WithOptions(t *testing.T) {
	mockService := &MockLanguageDetectionService{}

	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(60 * time.Second),
	}

	server := NewServer(mockService, opts...)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}
}

func TestServer_DetectLanguage_Success(t *testing.T) {
	ctx := context.Background()

	mockResponse := &domain.LanguageDetectionResponse{
		LanguageCode: "en-US",
		Confidence:   0.95,
		Alternatives: []domain.LanguageAlternative{
			{
				LanguageCode: "es-ES",
				Confidence:   0.85,
			},
		},
		DocumentID: "doc-123",
		Metadata: domain.ProcessingMetadata{
			ProcessingTimeMs: 100,
			ServiceVersion:   "1.0.0",
			ModelVersion:     "1.0.0",
			Provider:         "test",
			Details: map[string]string{
				"test": "value",
			},
		},
	}

	mockService := &MockLanguageDetectionService{
		response: mockResponse,
	}

	server := NewServer(mockService)

	req := &pb.DetectLanguageRequest{
		Text:       "Hello world",
		DocumentId: "doc-123",
		Metadata: map[string]string{
			"source": "test",
		},
	}

	resp, err := server.DetectLanguage(ctx, req)

	if err != nil {
		t.Fatalf("DetectLanguage() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.LanguageCode != "en-US" {
		t.Errorf("Expected language code 'en-US', got %s", resp.LanguageCode)
	}

	if resp.Confidence != 0.95 {
		t.Errorf("Expected confidence 0.95, got %f", resp.Confidence)
	}

	if resp.DocumentId != "doc-123" {
		t.Errorf("Expected document ID 'doc-123', got %s", resp.DocumentId)
	}

	if len(resp.Alternatives) != 1 {
		t.Errorf("Expected 1 alternative, got %d", len(resp.Alternatives))
	}

	if resp.Alternatives[0].LanguageCode != "es-ES" {
		t.Errorf("Expected alternative language 'es-ES', got %s", resp.Alternatives[0].LanguageCode)
	}

	if resp.Alternatives[0].Confidence != 0.85 {
		t.Errorf("Expected alternative confidence 0.85, got %f", resp.Alternatives[0].Confidence)
	}

	if resp.Metadata == nil {
		t.Fatal("Expected metadata, got nil")
	}

	if resp.Metadata.ProcessingTimeMs != 100 {
		t.Errorf("Expected processing time 100ms, got %d", resp.Metadata.ProcessingTimeMs)
	}

	if resp.Metadata.ServiceVersion != "1.0.0" {
		t.Errorf("Expected service version '1.0.0', got %s", resp.Metadata.ServiceVersion)
	}

	if resp.Metadata.ModelVersion != "1.0.0" {
		t.Errorf("Expected model version '1.0.0', got %s", resp.Metadata.ModelVersion)
	}

	if resp.Metadata.Provider != "test" {
		t.Errorf("Expected provider 'test', got %s", resp.Metadata.Provider)
	}
}

func TestServer_DetectLanguage_ServiceError(t *testing.T) {
	ctx := context.Background()

	mockService := &MockLanguageDetectionService{
		err: domain.ErrEmptyText,
	}

	server := NewServer(mockService)

	req := &pb.DetectLanguageRequest{
		Text: "",
	}

	_, err := server.DetectLanguage(ctx, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error is properly wrapped
	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}

func TestServer_DetectLanguage_EmptyRequest(t *testing.T) {
	ctx := context.Background()

	mockService := &MockLanguageDetectionService{
		err: domain.ErrEmptyText,
	}
	server := NewServer(mockService)

	req := &pb.DetectLanguageRequest{
		Text: "",
	}

	_, err := server.DetectLanguage(ctx, req)

	if err == nil {
		t.Fatal("Expected error for empty request, got nil")
	}
}

func TestServer_ConvertToProtobufResponse(t *testing.T) {
	server := &Server{}

	domainResp := &domain.LanguageDetectionResponse{
		LanguageCode: "fr-FR",
		Confidence:   0.87,
		Alternatives: []domain.LanguageAlternative{
			{
				LanguageCode: "de-DE",
				Confidence:   0.75,
			},
			{
				LanguageCode: "it-IT",
				Confidence:   0.65,
			},
		},
		DocumentID: "doc-456",
		Metadata: domain.ProcessingMetadata{
			ProcessingTimeMs: 250,
			ServiceVersion:   "2.0.0",
			ModelVersion:     "2.0.0",
			Provider:         "aws-comprehend",
		},
	}

	pbResp := server.convertToProtobufResponse(domainResp)

	if pbResp == nil {
		t.Fatal("Expected protobuf response, got nil")
	}

	if pbResp.LanguageCode != "fr-FR" {
		t.Errorf("Expected language code 'fr-FR', got %s", pbResp.LanguageCode)
	}

	if pbResp.Confidence != 0.87 {
		t.Errorf("Expected confidence 0.87, got %f", pbResp.Confidence)
	}

	if pbResp.DocumentId != "doc-456" {
		t.Errorf("Expected document ID 'doc-456', got %s", pbResp.DocumentId)
	}

	if len(pbResp.Alternatives) != 2 {
		t.Errorf("Expected 2 alternatives, got %d", len(pbResp.Alternatives))
	}

	if pbResp.Alternatives[0].LanguageCode != "de-DE" {
		t.Errorf("Expected first alternative 'de-DE', got %s", pbResp.Alternatives[0].LanguageCode)
	}

	if pbResp.Alternatives[0].Confidence != 0.75 {
		t.Errorf("Expected first alternative confidence 0.75, got %f", pbResp.Alternatives[0].Confidence)
	}

	if pbResp.Alternatives[1].LanguageCode != "it-IT" {
		t.Errorf("Expected second alternative 'it-IT', got %s", pbResp.Alternatives[1].LanguageCode)
	}

	if pbResp.Alternatives[1].Confidence != 0.65 {
		t.Errorf("Expected second alternative confidence 0.65, got %f", pbResp.Alternatives[1].Confidence)
	}

	if pbResp.Metadata == nil {
		t.Fatal("Expected metadata, got nil")
	}

	if pbResp.Metadata.ProcessingTimeMs != 250 {
		t.Errorf("Expected processing time 250ms, got %d", pbResp.Metadata.ProcessingTimeMs)
	}

	if pbResp.Metadata.ServiceVersion != "2.0.0" {
		t.Errorf("Expected service version '2.0.0', got %s", pbResp.Metadata.ServiceVersion)
	}

	if pbResp.Metadata.ModelVersion != "2.0.0" {
		t.Errorf("Expected model version '2.0.0', got %s", pbResp.Metadata.ModelVersion)
	}

	if pbResp.Metadata.Provider != "aws-comprehend" {
		t.Errorf("Expected provider 'aws-comprehend', got %s", pbResp.Metadata.Provider)
	}
}

func TestServer_ConvertToProtobufResponse_NoAlternatives(t *testing.T) {
	server := &Server{}

	domainResp := &domain.LanguageDetectionResponse{
		LanguageCode: "ja-JP",
		Confidence:   0.92,
		Alternatives: []domain.LanguageAlternative{}, // Empty alternatives
		DocumentID:   "doc-789",
		Metadata: domain.ProcessingMetadata{
			ProcessingTimeMs: 150,
			ServiceVersion:   "1.0.0",
			ModelVersion:     "1.0.0",
			Provider:         "fallback",
		},
	}

	pbResp := server.convertToProtobufResponse(domainResp)

	if pbResp == nil {
		t.Fatal("Expected protobuf response, got nil")
	}

	if pbResp.LanguageCode != "ja-JP" {
		t.Errorf("Expected language code 'ja-JP', got %s", pbResp.LanguageCode)
	}

	if len(pbResp.Alternatives) != 0 {
		t.Errorf("Expected 0 alternatives, got %d", len(pbResp.Alternatives))
	}
}

func TestServer_ConvertToProtobufResponse_NilAlternatives(t *testing.T) {
	server := &Server{}

	domainResp := &domain.LanguageDetectionResponse{
		LanguageCode: "ko-KR",
		Confidence:   0.88,
		Alternatives: nil, // Nil alternatives
		DocumentID:   "doc-101",
		Metadata: domain.ProcessingMetadata{
			ProcessingTimeMs: 200,
			ServiceVersion:   "1.0.0",
			ModelVersion:     "1.0.0",
			Provider:         "aws-comprehend",
		},
	}

	pbResp := server.convertToProtobufResponse(domainResp)

	if pbResp == nil {
		t.Fatal("Expected protobuf response, got nil")
	}

	if pbResp.LanguageCode != "ko-KR" {
		t.Errorf("Expected language code 'ko-KR', got %s", pbResp.LanguageCode)
	}

	if len(pbResp.Alternatives) != 0 {
		t.Errorf("Expected 0 alternatives, got %d", len(pbResp.Alternatives))
	}
}

func TestServer_StartWithContext_InvalidAddress(t *testing.T) {
	mockService := &MockLanguageDetectionService{}
	server := NewServer(mockService)

	ctx := context.Background()
	invalidAddress := "invalid:address:format"

	err := server.StartWithContext(ctx, invalidAddress)

	if err == nil {
		t.Fatal("Expected error for invalid address, got nil")
	}
}

func TestServer_StartWithContext_CancelledContext(t *testing.T) {
	mockService := &MockLanguageDetectionService{}
	server := NewServer(mockService)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := server.StartWithContext(ctx, "localhost:0")

	if err == nil {
		t.Fatal("Expected error for cancelled context, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestServer_Stop(t *testing.T) {
	mockService := &MockLanguageDetectionService{}
	server := NewServer(mockService)

	// Test stopping a server that was never started
	err := server.Stop()

	if err != nil {
		t.Errorf("Stop() error = %v, want nil", err)
	}
}

func TestServer_HealthCheck(t *testing.T) {
	mockService := &MockLanguageDetectionService{}
	server := NewServer(mockService)

	// The health server should be initialized and set to serving
	if server.healthServer == nil {
		t.Fatal("Expected health server to be initialized")
	}

	// Test health check
	ctx := context.Background()
	healthReq := &grpc_health_v1.HealthCheckRequest{
		Service: "language_detection.LanguageDetectionService",
	}

	resp, err := server.healthServer.Check(ctx, healthReq)

	if err != nil {
		t.Errorf("Health check error = %v, want nil", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("Expected health status SERVING, got %v", resp.Status)
	}
}

func TestServer_Integration(t *testing.T) {
	ctx := context.Background()

	// Create a realistic response
	mockResponse := &domain.LanguageDetectionResponse{
		LanguageCode: "zh-CN",
		Confidence:   0.93,
		Alternatives: []domain.LanguageAlternative{
			{
				LanguageCode: "zh-TW",
				Confidence:   0.78,
			},
		},
		DocumentID: "doc-integration",
		Metadata: domain.ProcessingMetadata{
			ProcessingTimeMs: 300,
			ServiceVersion:   "1.0.0",
			ModelVersion:     "1.0.0",
			Provider:         "aws-comprehend",
			Details: map[string]string{
				"region":        "us-east-1",
				"total_langs":   "3",
				"aws_lang_code": "zh",
			},
		},
	}

	mockService := &MockLanguageDetectionService{
		response: mockResponse,
	}

	server := NewServer(mockService)

	req := &pb.DetectLanguageRequest{
		Text:       "你好世界",
		DocumentId: "doc-integration",
		Metadata: map[string]string{
			"source": "integration-test",
		},
	}

	resp, err := server.DetectLanguage(ctx, req)

	if err != nil {
		t.Fatalf("Integration test error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("Expected integration test response, got nil")
	}

	// Verify all fields are properly converted
	if resp.LanguageCode != "zh-CN" {
		t.Errorf("Integration test: expected language 'zh-CN', got %s", resp.LanguageCode)
	}

	if resp.Confidence != 0.93 {
		t.Errorf("Integration test: expected confidence 0.93, got %f", resp.Confidence)
	}

	if resp.DocumentId != "doc-integration" {
		t.Errorf("Integration test: expected document ID 'doc-integration', got %s", resp.DocumentId)
	}

	if len(resp.Alternatives) != 1 {
		t.Errorf("Integration test: expected 1 alternative, got %d", len(resp.Alternatives))
	}

	if resp.Metadata == nil {
		t.Fatal("Integration test: expected metadata, got nil")
	}

	if resp.Metadata.ProcessingTimeMs != 300 {
		t.Errorf("Integration test: expected processing time 300ms, got %d", resp.Metadata.ProcessingTimeMs)
	}

	if resp.Metadata.Provider != "aws-comprehend" {
		t.Errorf("Integration test: expected provider 'aws-comprehend', got %s", resp.Metadata.Provider)
	}
}
