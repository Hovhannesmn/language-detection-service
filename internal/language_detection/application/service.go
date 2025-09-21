package application

import (
	"context"
	"fmt"
	"time"

	"language-detection-service/internal/language_detection/domain"
)

// LanguageDetectionServiceImpl implements the LanguageDetectionService interface
type LanguageDetectionServiceImpl struct {
	detector domain.LanguageDetector
	config   domain.ConfigProvider
}

// NewLanguageDetectionService creates a new language detection service
func NewLanguageDetectionService(
	detector domain.LanguageDetector,
	config domain.ConfigProvider,
) *LanguageDetectionServiceImpl {
	return &LanguageDetectionServiceImpl{
		detector: detector,
		config:   config,
	}
}

// DetectLanguage performs language detection with business logic
func (s *LanguageDetectionServiceImpl) DetectLanguage(
	ctx context.Context,
	request *domain.LanguageDetectionRequest,
) (*domain.LanguageDetectionResponse, error) {
	startTime := time.Now()

	// Validate input
	if err := s.validateRequest(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Perform language detection
	response, err := s.detector.DetectLanguage(ctx, request.Text)
	if err != nil {
		return nil, fmt.Errorf("language detection failed: %w", err)
	}

	// Validate response
	if err := s.validateResponse(response); err != nil {
		return nil, fmt.Errorf("response validation failed: %w", err)
	}

	// Update metadata
	response.DocumentID = request.DocumentID
	response.Metadata.ProcessingTimeMs = time.Since(startTime).Milliseconds()
	response.Metadata.ServiceVersion = s.config.GetServiceVersion()
	response.Metadata.ModelVersion = s.config.GetModelVersion()

	return response, nil
}

// validateRequest validates the incoming request
func (s *LanguageDetectionServiceImpl) validateRequest(request *domain.LanguageDetectionRequest) error {
	if request == nil {
		return domain.ErrInvalidRequest
	}

	text := string(request.Text)
	if text == "" {
		return domain.ErrEmptyText
	}

	if len(text) > s.config.GetMaxTextLength() {
		return fmt.Errorf("%w: text length %d exceeds maximum %d",
			domain.ErrTextTooLong, len(text), s.config.GetMaxTextLength())
	}

	return nil
}

// validateResponse validates the detection response
func (s *LanguageDetectionServiceImpl) validateResponse(response *domain.LanguageDetectionResponse) error {
	if response == nil {
		return domain.ErrInternalError
	}

	if float32(response.Confidence) < s.config.GetMinConfidenceThreshold() {
		return fmt.Errorf("%w: confidence %.2f below threshold %.2f",
			domain.ErrLowConfidence, float32(response.Confidence), s.config.GetMinConfidenceThreshold())
	}

	// Check if language is supported
	supported := s.config.GetSupportedLanguages()
	if len(supported) > 0 {
		found := false
		for _, lang := range supported {
			if lang == response.LanguageCode {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%w: %s", domain.ErrInvalidLanguageCode, response.LanguageCode)
		}
	}

	return nil
}
