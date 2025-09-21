package domain

import "context"

// LanguageDetector defines the port for language detection services
type LanguageDetector interface {
	// DetectLanguage detects the dominant language of the given text
	DetectLanguage(ctx context.Context, text Text) (*LanguageDetectionResponse, error)
}

// LanguageDetectionService defines the application service port
type LanguageDetectionService interface {
	// DetectLanguage performs language detection with business logic
	DetectLanguage(ctx context.Context, request *LanguageDetectionRequest) (*LanguageDetectionResponse, error)
}

// ConfigProvider defines the port for configuration
type ConfigProvider interface {
	// GetMaxTextLength returns the maximum allowed text length
	GetMaxTextLength() int
	
	// GetMinConfidenceThreshold returns the minimum confidence threshold
	GetMinConfidenceThreshold() float32
	
	// GetSupportedLanguages returns the list of supported languages
	GetSupportedLanguages() []LanguageCode
	
	// GetServiceVersion returns the service version
	GetServiceVersion() string
	
	// GetModelVersion returns the model version
	GetModelVersion() string
}
