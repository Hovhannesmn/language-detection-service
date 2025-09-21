package domain

// LanguageCode represents a language code (e.g., "en-US", "es-ES")
type LanguageCode string

// Confidence represents the confidence score for language detection (0.0 to 1.0)
type Confidence float32

// Text represents the input text for language detection
type Text string

// LanguageDetectionRequest represents a request for language detection
type LanguageDetectionRequest struct {
	Text       Text              `json:"text"`
	DocumentID string            `json:"document_id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// LanguageDetectionResponse represents the response from language detection
type LanguageDetectionResponse struct {
	LanguageCode LanguageCode           `json:"language_code"`
	Confidence   Confidence             `json:"confidence"`
	Alternatives []LanguageAlternative  `json:"alternatives,omitempty"`
	DocumentID   string                 `json:"document_id,omitempty"`
	Metadata     ProcessingMetadata     `json:"metadata"`
}

// LanguageAlternative represents an alternative language detection result
type LanguageAlternative struct {
	LanguageCode LanguageCode `json:"language_code"`
	Confidence   Confidence   `json:"confidence"`
}

// ProcessingMetadata contains information about the processing
type ProcessingMetadata struct {
	ProcessingTimeMs int64             `json:"processing_time_ms"`
	ServiceVersion   string            `json:"service_version"`
	ModelVersion     string            `json:"model_version"`
	Provider         string            `json:"provider"`
	Details          map[string]string `json:"details,omitempty"`
}
