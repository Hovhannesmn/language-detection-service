package application

import (
	"context"
	"errors"
	"testing"

	"language-detection-service/internal/language_detection/domain"
)

// MockLanguageDetector is a mock implementation of LanguageDetector
type MockLanguageDetector struct {
	response *domain.LanguageDetectionResponse
	err      error
}

func (m *MockLanguageDetector) DetectLanguage(ctx context.Context, text domain.Text) (*domain.LanguageDetectionResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

// MockConfigProvider is a mock implementation of ConfigProvider
type MockConfigProvider struct {
	maxTextLength          int
	minConfidenceThreshold float32
	supportedLanguages     []domain.LanguageCode
	serviceVersion         string
	modelVersion           string
}

func (m *MockConfigProvider) GetMaxTextLength() int {
	return m.maxTextLength
}

func (m *MockConfigProvider) GetMinConfidenceThreshold() float32 {
	return m.minConfidenceThreshold
}

func (m *MockConfigProvider) GetSupportedLanguages() []domain.LanguageCode {
	return m.supportedLanguages
}

func (m *MockConfigProvider) GetServiceVersion() string {
	return m.serviceVersion
}

func (m *MockConfigProvider) GetModelVersion() string {
	return m.modelVersion
}

func TestNewLanguageDetectionService(t *testing.T) {
	detector := &MockLanguageDetector{}
	config := &MockConfigProvider{}
	
	service := NewLanguageDetectionService(detector, config)
	
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
	
	if service.detector != detector {
		t.Errorf("Expected detector to be %v, got %v", detector, service.detector)
	}
	
	if service.config != config {
		t.Errorf("Expected config to be %v, got %v", config, service.config)
	}
}

func TestDetectLanguage_Success(t *testing.T) {
	ctx := context.Background()
	
	mockResponse := &domain.LanguageDetectionResponse{
		LanguageCode: "en-US",
		Confidence:   0.95,
		Metadata: domain.ProcessingMetadata{
			Provider: "test",
		},
	}
	
	detector := &MockLanguageDetector{
		response: mockResponse,
	}
	
	config := &MockConfigProvider{
		maxTextLength:          1000,
		minConfidenceThreshold: 0.1,
		supportedLanguages:     []domain.LanguageCode{"en-US", "es-ES"},
		serviceVersion:         "1.0.0",
		modelVersion:           "1.0.0",
	}
	
	service := NewLanguageDetectionService(detector, config)
	
	request := &domain.LanguageDetectionRequest{
		Text:       "Hello world",
		DocumentID: "doc-123",
	}
	
	response, err := service.DetectLanguage(ctx, request)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if response.LanguageCode != "en-US" {
		t.Errorf("Expected language code 'en-US', got %v", response.LanguageCode)
	}
	
	if response.DocumentID != "doc-123" {
		t.Errorf("Expected document ID 'doc-123', got %v", response.DocumentID)
	}
	
	if response.Metadata.ServiceVersion != "1.0.0" {
		t.Errorf("Expected service version '1.0.0', got %v", response.Metadata.ServiceVersion)
	}
	
	if response.Metadata.ModelVersion != "1.0.0" {
		t.Errorf("Expected model version '1.0.0', got %v", response.Metadata.ModelVersion)
	}
}

func TestDetectLanguage_DetectorError(t *testing.T) {
	ctx := context.Background()
	
	detector := &MockLanguageDetector{
		err: errors.New("detector error"),
	}
	
	config := &MockConfigProvider{
		maxTextLength:          1000,
		minConfidenceThreshold: 0.1,
	}
	
	service := NewLanguageDetectionService(detector, config)
	
	request := &domain.LanguageDetectionRequest{
		Text: "Hello world",
	}
	
	_, err := service.DetectLanguage(ctx, request)
	
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetectLanguage_EmptyText(t *testing.T) {
	ctx := context.Background()
	
	detector := &MockLanguageDetector{}
	config := &MockConfigProvider{
		maxTextLength: 1000,
	}
	
	service := NewLanguageDetectionService(detector, config)
	
	request := &domain.LanguageDetectionRequest{
		Text: "",
	}
	
	_, err := service.DetectLanguage(ctx, request)
	
	if err == nil {
		t.Fatal("Expected error for empty text, got nil")
	}
	
	if !errors.Is(err, domain.ErrEmptyText) {
		t.Errorf("Expected ErrEmptyText, got %v", err)
	}
}

func TestDetectLanguage_NilRequest(t *testing.T) {
	ctx := context.Background()
	
	detector := &MockLanguageDetector{}
	config := &MockConfigProvider{}
	
	service := NewLanguageDetectionService(detector, config)
	
	_, err := service.DetectLanguage(ctx, nil)
	
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	
	if !errors.Is(err, domain.ErrInvalidRequest) {
		t.Errorf("Expected ErrInvalidRequest, got %v", err)
	}
}

func TestDetectLanguage_TextTooLong(t *testing.T) {
	ctx := context.Background()
	
	detector := &MockLanguageDetector{}
	config := &MockConfigProvider{
		maxTextLength: 10,
	}
	
	service := NewLanguageDetectionService(detector, config)
	
	request := &domain.LanguageDetectionRequest{
		Text: "This text is too long",
	}
	
	_, err := service.DetectLanguage(ctx, request)
	
	if err == nil {
		t.Fatal("Expected error for text too long, got nil")
	}
	
	if !errors.Is(err, domain.ErrTextTooLong) {
		t.Errorf("Expected ErrTextTooLong, got %v", err)
	}
}

func TestDetectLanguage_LowConfidence(t *testing.T) {
	ctx := context.Background()
	
	mockResponse := &domain.LanguageDetectionResponse{
		LanguageCode: "en-US",
		Confidence:   0.05, // Below threshold
		Metadata: domain.ProcessingMetadata{
			Provider: "test",
		},
	}
	
	detector := &MockLanguageDetector{
		response: mockResponse,
	}
	
	config := &MockConfigProvider{
		maxTextLength:          1000,
		minConfidenceThreshold: 0.1,
	}
	
	service := NewLanguageDetectionService(detector, config)
	
	request := &domain.LanguageDetectionRequest{
		Text: "Hello world",
	}
	
	_, err := service.DetectLanguage(ctx, request)
	
	if err == nil {
		t.Fatal("Expected error for low confidence, got nil")
	}
	
	if !errors.Is(err, domain.ErrLowConfidence) {
		t.Errorf("Expected ErrLowConfidence, got %v", err)
	}
}

func TestDetectLanguage_UnsupportedLanguage(t *testing.T) {
	ctx := context.Background()
	
	mockResponse := &domain.LanguageDetectionResponse{
		LanguageCode: "fr-FR",
		Confidence:   0.95,
		Metadata: domain.ProcessingMetadata{
			Provider: "test",
		},
	}
	
	detector := &MockLanguageDetector{
		response: mockResponse,
	}
	
	config := &MockConfigProvider{
		maxTextLength:          1000,
		minConfidenceThreshold: 0.1,
		supportedLanguages:     []domain.LanguageCode{"en-US", "es-ES"}, // fr-FR not supported
	}
	
	service := NewLanguageDetectionService(detector, config)
	
	request := &domain.LanguageDetectionRequest{
		Text: "Bonjour monde",
	}
	
	_, err := service.DetectLanguage(ctx, request)
	
	if err == nil {
		t.Fatal("Expected error for unsupported language, got nil")
	}
	
	if !errors.Is(err, domain.ErrInvalidLanguageCode) {
		t.Errorf("Expected ErrInvalidLanguageCode, got %v", err)
	}
}

func TestDetectLanguage_ProcessingTime(t *testing.T) {
	ctx := context.Background()
	
	mockResponse := &domain.LanguageDetectionResponse{
		LanguageCode: "en-US",
		Confidence:   0.95,
		Metadata: domain.ProcessingMetadata{
			Provider: "test",
		},
	}
	
	detector := &MockLanguageDetector{
		response: mockResponse,
	}
	
	config := &MockConfigProvider{
		maxTextLength:          1000,
		minConfidenceThreshold: 0.1,
		serviceVersion:         "1.0.0",
		modelVersion:           "1.0.0",
	}
	
	service := NewLanguageDetectionService(detector, config)
	
	request := &domain.LanguageDetectionRequest{
		Text: "Hello world",
	}
	
	response, err := service.DetectLanguage(ctx, request)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Processing time should be >= 0 (could be 0 for very fast tests)
	if response.Metadata.ProcessingTimeMs < 0 {
		t.Errorf("Expected processing time >= 0, got %d", response.Metadata.ProcessingTimeMs)
	}
}

func TestValidateRequest(t *testing.T) {
	service := &LanguageDetectionServiceImpl{
		config: &MockConfigProvider{
			maxTextLength: 100,
		},
	}
	
	tests := []struct {
		name    string
		request *domain.LanguageDetectionRequest
		wantErr error
	}{
		{
			name:    "nil request",
			request: nil,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "empty text",
			request: &domain.LanguageDetectionRequest{
				Text: "",
			},
			wantErr: domain.ErrEmptyText,
		},
		{
			name: "text too long",
			request: &domain.LanguageDetectionRequest{
				Text: "This text is way too long and exceeds the maximum allowed length for testing purposes and continues to be very long indeed with many more characters added to ensure it definitely exceeds the limit",
			},
			wantErr: domain.ErrTextTooLong,
		},
		{
			name: "valid request",
			request: &domain.LanguageDetectionRequest{
				Text: "Hello world",
			},
			wantErr: nil,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateRequest(tt.request)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestValidateResponse(t *testing.T) {
	service := &LanguageDetectionServiceImpl{
		config: &MockConfigProvider{
			minConfidenceThreshold: 0.1,
			supportedLanguages:     []domain.LanguageCode{"en-US", "es-ES"},
		},
	}
	
	tests := []struct {
		name     string
		response *domain.LanguageDetectionResponse
		wantErr  error
	}{
		{
			name:     "nil response",
			response: nil,
			wantErr:  domain.ErrInternalError,
		},
		{
			name: "low confidence",
			response: &domain.LanguageDetectionResponse{
				LanguageCode: "en-US",
				Confidence:   0.05,
			},
			wantErr: domain.ErrLowConfidence,
		},
		{
			name: "unsupported language",
			response: &domain.LanguageDetectionResponse{
				LanguageCode: "fr-FR",
				Confidence:   0.95,
			},
			wantErr: domain.ErrInvalidLanguageCode,
		},
		{
			name: "valid response",
			response: &domain.LanguageDetectionResponse{
				LanguageCode: "en-US",
				Confidence:   0.95,
			},
			wantErr: nil,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateResponse(tt.response)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("validateResponse() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("validateResponse() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
