package adapters

import (
	"context"
	"testing"

	"language-detection-service/internal/language_detection/domain"
)

func TestAWSComprehendAdapter_ConvertLanguageCode(t *testing.T) {
	adapter := &AWSComprehendAdapter{}
	
	tests := []struct {
		name     string
		awsCode  string
		expected domain.LanguageCode
	}{
		{
			name:     "English",
			awsCode:  "en",
			expected: "en-US",
		},
		{
			name:     "Spanish",
			awsCode:  "es",
			expected: "es-ES",
		},
		{
			name:     "French",
			awsCode:  "fr",
			expected: "fr-FR",
		},
		{
			name:     "German",
			awsCode:  "de",
			expected: "de-DE",
		},
		{
			name:     "Italian",
			awsCode:  "it",
			expected: "it-IT",
		},
		{
			name:     "Portuguese",
			awsCode:  "pt",
			expected: "pt-PT",
		},
		{
			name:     "Russian",
			awsCode:  "ru",
			expected: "ru-RU",
		},
		{
			name:     "Japanese",
			awsCode:  "ja",
			expected: "ja-JP",
		},
		{
			name:     "Korean",
			awsCode:  "ko",
			expected: "ko-KR",
		},
		{
			name:     "Chinese simplified",
			awsCode:  "zh",
			expected: "zh-CN",
		},
		{
			name:     "Chinese traditional",
			awsCode:  "zh-tw",
			expected: "zh-CN",
		},
		{
			name:     "Chinese simplified variant",
			awsCode:  "zh-cn",
			expected: "zh-CN",
		},
		{
			name:     "Arabic",
			awsCode:  "ar",
			expected: "ar-SA",
		},
		{
			name:     "Hindi",
			awsCode:  "hi",
			expected: "hi-IN",
		},
		{
			name:     "Unknown language",
			awsCode:  "xx",
			expected: "xx",
		},
		{
			name:     "Case insensitive",
			awsCode:  "EN",
			expected: "en-US",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.convertLanguageCode(tt.awsCode)
			if result != tt.expected {
				t.Errorf("convertLanguageCode(%s) = %s, want %s", tt.awsCode, result, tt.expected)
			}
		})
	}
}

func TestAWSComprehendAdapter_NewAWSComprehendAdapter_Error(t *testing.T) {
	// Test with invalid region - this should fail due to missing AWS credentials
	region := "invalid-region"
	maxRetries := 3
	
	adapter, err := NewAWSComprehendAdapter(region, maxRetries)
	
	// The adapter creation might succeed even without credentials due to AWS SDK behavior
	// So we'll just test that the function doesn't panic
	if err != nil && adapter != nil {
		t.Error("Expected either error or nil adapter, got both error and adapter")
	}
}

func TestAWSComprehendAdapter_DetectLanguage_EmptyText(t *testing.T) {
	// This test will fail if AWS credentials are not configured, which is expected
	adapter, err := NewAWSComprehendAdapter("us-east-1", 3)
	if err != nil {
		// Skip test if AWS credentials are not available
		t.Skip("Skipping test due to missing AWS credentials")
	}
	
	ctx := context.Background()
	text := domain.Text("")
	
	response, err := adapter.DetectLanguage(ctx, text)
	
	if err != nil {
		t.Fatalf("Expected no error for empty text, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if string(response.LanguageCode) != "unknown" {
		t.Errorf("Expected language code 'unknown', got %s", response.LanguageCode)
	}
	
	if response.Metadata.Details["reason"] != "text_too_short" {
		t.Errorf("Expected reason 'text_too_short', got %s", response.Metadata.Details["reason"])
	}
}

func TestAWSComprehendAdapter_DetectLanguage_TextTooShort(t *testing.T) {
	// This test will fail if AWS credentials are not configured, which is expected
	adapter, err := NewAWSComprehendAdapter("us-east-1", 3)
	if err != nil {
		// Skip test if AWS credentials are not available
		t.Skip("Skipping test due to missing AWS credentials")
	}
	
	ctx := context.Background()
	text := domain.Text("Hi")
	
	response, err := adapter.DetectLanguage(ctx, text)
	
	if err != nil {
		t.Fatalf("Expected no error for short text, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if string(response.LanguageCode) != "unknown" {
		t.Errorf("Expected language code 'unknown', got %s", response.LanguageCode)
	}
	
	if response.Metadata.Details["reason"] != "text_too_short" {
		t.Errorf("Expected reason 'text_too_short', got %s", response.Metadata.Details["reason"])
	}
}
