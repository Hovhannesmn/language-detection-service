package domain

import (
	"testing"
)

func TestLanguageCode(t *testing.T) {
	tests := []struct {
		name     string
		code     LanguageCode
		expected string
	}{
		{
			name:     "English US",
			code:     "en-US",
			expected: "en-US",
		},
		{
			name:     "Spanish Spain",
			code:     "es-ES",
			expected: "es-ES",
		},
		{
			name:     "Unknown language",
			code:     "unknown",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.code) != tt.expected {
				t.Errorf("LanguageCode = %v, want %v", tt.code, tt.expected)
			}
		})
	}
}

func TestConfidence(t *testing.T) {
	tests := []struct {
		name     string
		conf     Confidence
		expected float32
	}{
		{
			name:     "High confidence",
			conf:     0.95,
			expected: 0.95,
		},
		{
			name:     "Low confidence",
			conf:     0.1,
			expected: 0.1,
		},
		{
			name:     "Zero confidence",
			conf:     0.0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if float32(tt.conf) != tt.expected {
				t.Errorf("Confidence = %v, want %v", tt.conf, tt.expected)
			}
		})
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		name     string
		text     Text
		expected string
	}{
		{
			name:     "English text",
			text:     "Hello world",
			expected: "Hello world",
		},
		{
			name:     "Spanish text",
			text:     "Hola mundo",
			expected: "Hola mundo",
		},
		{
			name:     "Empty text",
			text:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.text) != tt.expected {
				t.Errorf("Text = %v, want %v", tt.text, tt.expected)
			}
		})
	}
}

func TestLanguageDetectionRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  LanguageDetectionRequest
		expected string
	}{
		{
			name: "Valid request",
			request: LanguageDetectionRequest{
				Text:       "Hello world",
				DocumentID: "doc-123",
				Metadata: map[string]string{
					"source": "test",
				},
			},
			expected: "Hello world",
		},
		{
			name: "Request without document ID",
			request: LanguageDetectionRequest{
				Text: "Hola mundo",
			},
			expected: "Hola mundo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.request.Text) != tt.expected {
				t.Errorf("LanguageDetectionRequest.Text = %v, want %v", tt.request.Text, tt.expected)
			}
		})
	}
}

func TestLanguageDetectionResponse(t *testing.T) {
	tests := []struct {
		name     string
		response LanguageDetectionResponse
		expected string
	}{
		{
			name: "Valid response",
			response: LanguageDetectionResponse{
				LanguageCode: "en-US",
				Confidence:   0.95,
				DocumentID:   "doc-123",
				Metadata: ProcessingMetadata{
					ProcessingTimeMs: 100,
					ServiceVersion:   "1.0.0",
					ModelVersion:     "1.0.0",
					Provider:         "test",
				},
			},
			expected: "en-US",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.response.LanguageCode) != tt.expected {
				t.Errorf("LanguageDetectionResponse.LanguageCode = %v, want %v", tt.response.LanguageCode, tt.expected)
			}
		})
	}
}

func TestLanguageAlternative(t *testing.T) {
	tests := []struct {
		name     string
		alt      LanguageAlternative
		expected string
	}{
		{
			name: "Valid alternative",
			alt: LanguageAlternative{
				LanguageCode: "es-ES",
				Confidence:   0.85,
			},
			expected: "es-ES",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.alt.LanguageCode) != tt.expected {
				t.Errorf("LanguageAlternative.LanguageCode = %v, want %v", tt.alt.LanguageCode, tt.expected)
			}
		})
	}
}

func TestProcessingMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata ProcessingMetadata
		expected string
	}{
		{
			name: "Valid metadata",
			metadata: ProcessingMetadata{
				ProcessingTimeMs: 150,
				ServiceVersion:   "1.0.0",
				ModelVersion:     "1.0.0",
				Provider:         "aws-comprehend",
				Details: map[string]string{
					"region": "us-east-1",
				},
			},
			expected: "aws-comprehend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.metadata.Provider != tt.expected {
				t.Errorf("ProcessingMetadata.Provider = %v, want %v", tt.metadata.Provider, tt.expected)
			}
		})
	}
}
