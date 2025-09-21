package adapters

import (
	"context"
	"testing"

	"language-detection-service/internal/language_detection/domain"
)

func TestNewFallbackAdapter(t *testing.T) {
	adapter := NewFallbackAdapter()
	
	if adapter == nil {
		t.Fatal("Expected adapter to be created, got nil")
	}
	
	if adapter.patterns == nil {
		t.Fatal("Expected patterns to be initialized, got nil")
	}
	
	// Check that we have patterns for expected languages
	expectedLanguages := []string{"en-US", "es-ES", "fr-FR", "de-DE"}
	for _, lang := range expectedLanguages {
		if _, exists := adapter.patterns[lang]; !exists {
			t.Errorf("Expected patterns for language %s, but not found", lang)
		}
	}
}

func TestFallbackAdapter_DetectLanguage_English(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	tests := []struct {
		name     string
		text     string
		expected string
		minConf  float32
	}{
		{
			name:     "Simple English text",
			text:     "Hello world",
			expected: "en-US",
			minConf:  0.15,
		},
		{
			name:     "English with common words",
			text:     "The quick brown fox jumps over the lazy dog",
			expected: "en-US",
			minConf:  0.15,
		},
		{
			name:     "English presentation text",
			text:     "Welcome to today's presentation about technology trends",
			expected: "en-US",
			minConf:  0.15,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := adapter.DetectLanguage(ctx, domain.Text(tt.text))
			
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			
			if response == nil {
				t.Fatal("Expected response, got nil")
			}
			
			if string(response.LanguageCode) != tt.expected {
				t.Errorf("Expected language code %s, got %s", tt.expected, response.LanguageCode)
			}
			
			if float32(response.Confidence) < tt.minConf {
				t.Errorf("Expected confidence >= %.2f, got %.2f", tt.minConf, response.Confidence)
			}
			
			if response.Metadata.Provider != "fallback" {
				t.Errorf("Expected provider 'fallback', got %s", response.Metadata.Provider)
			}
		})
	}
}

func TestFallbackAdapter_DetectLanguage_Spanish(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	tests := []struct {
		name     string
		text     string
		expected string
		minConf  float32
	}{
		{
			name:     "Simple Spanish text",
			text:     "Hola mundo",
			expected: "es-ES",
			minConf:  0.15,
		},
		{
			name:     "Spanish with common words",
			text:     "El perro marrón salta sobre el gato perezoso",
			expected: "es-ES",
			minConf:  0.15,
		},
		{
			name:     "Spanish presentation text",
			text:     "Bienvenidos a la presentación de hoy sobre tendencias tecnológicas",
			expected: "es-ES",
			minConf:  0.15,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := adapter.DetectLanguage(ctx, domain.Text(tt.text))
			
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			
			if response == nil {
				t.Fatal("Expected response, got nil")
			}
			
			if string(response.LanguageCode) != tt.expected {
				t.Errorf("Expected language code %s, got %s", tt.expected, response.LanguageCode)
			}
			
			if float32(response.Confidence) < tt.minConf {
				t.Errorf("Expected confidence >= %.2f, got %.2f", tt.minConf, response.Confidence)
			}
		})
	}
}

func TestFallbackAdapter_DetectLanguage_French(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	text := "Bonjour, bienvenue à la présentation d'aujourd'hui"
	response, err := adapter.DetectLanguage(ctx, domain.Text(text))
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if string(response.LanguageCode) != "fr-FR" {
		t.Errorf("Expected language code 'fr-FR', got %s", response.LanguageCode)
	}
	
	if float32(response.Confidence) < 0.15 {
		t.Errorf("Expected confidence >= 0.15, got %.2f", response.Confidence)
	}
}

func TestFallbackAdapter_DetectLanguage_German(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	text := "Willkommen zur heutigen Präsentation über Technologietrends"
	response, err := adapter.DetectLanguage(ctx, domain.Text(text))
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if string(response.LanguageCode) != "de-DE" {
		t.Errorf("Expected language code 'de-DE', got %s", response.LanguageCode)
	}
	
	if float32(response.Confidence) < 0.15 {
		t.Errorf("Expected confidence >= 0.15, got %.2f", response.Confidence)
	}
}

func TestFallbackAdapter_DetectLanguage_TextTooShort(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	tests := []struct {
		name string
		text string
	}{
		{
			name: "Empty text",
			text: "",
		},
		{
			name: "Single character",
			text: "a",
		},
		{
			name: "Two characters",
			text: "ab",
		},
		{
			name: "Whitespace only",
			text: "   ",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := adapter.DetectLanguage(ctx, domain.Text(tt.text))
			
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			
			if response == nil {
				t.Fatal("Expected response, got nil")
			}
			
			if string(response.LanguageCode) != "unknown" {
				t.Errorf("Expected language code 'unknown', got %s", response.LanguageCode)
			}
			
			if response.Confidence != 0 {
				t.Errorf("Expected confidence 0, got %.2f", response.Confidence)
			}
			
			if response.Metadata.Details["reason"] != "text_too_short" {
				t.Errorf("Expected reason 'text_too_short', got %s", response.Metadata.Details["reason"])
			}
		})
	}
}

func TestFallbackAdapter_DetectLanguage_NoWords(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	// Text with only punctuation and numbers
	text := "123 !@#$%^&*() 456"
	response, err := adapter.DetectLanguage(ctx, domain.Text(text))
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if string(response.LanguageCode) != "unknown" {
		t.Errorf("Expected language code 'unknown', got %s", response.LanguageCode)
	}
	
	// The reason could be either "no_words_found" or "low_confidence" depending on implementation
	reason := response.Metadata.Details["reason"]
	if reason != "no_words_found" && reason != "low_confidence" {
		t.Errorf("Expected reason 'no_words_found' or 'low_confidence', got %s", reason)
	}
}

func TestFallbackAdapter_DetectLanguage_LowConfidence(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	// Text with random words that don't match any patterns
	text := "xyzabc qwerty uiopas dfghjk"
	response, err := adapter.DetectLanguage(ctx, domain.Text(text))
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if string(response.LanguageCode) != "unknown" {
		t.Errorf("Expected language code 'unknown', got %s", response.LanguageCode)
	}
	
	if response.Metadata.Details["reason"] != "low_confidence" {
		t.Errorf("Expected reason 'low_confidence', got %s", response.Metadata.Details["reason"])
	}
}

func TestFallbackAdapter_DetectLanguage_Alternatives(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	// Mixed text that could be interpreted as multiple languages
	text := "Hello the world and hola mundo"
	response, err := adapter.DetectLanguage(ctx, domain.Text(text))
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	// Should have alternatives with confidence > 5%
	if len(response.Alternatives) == 0 {
		t.Error("Expected alternatives, got none")
	}
	
	for _, alt := range response.Alternatives {
		if float32(alt.Confidence) <= 0.05 {
			t.Errorf("Expected alternative confidence > 0.05, got %.2f", alt.Confidence)
		}
	}
}

func TestFallbackAdapter_DetectLanguage_Metadata(t *testing.T) {
	adapter := NewFallbackAdapter()
	ctx := context.Background()
	
	text := "Hello world"
	response, err := adapter.DetectLanguage(ctx, domain.Text(text))
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if response.Metadata.Provider != "fallback" {
		t.Errorf("Expected provider 'fallback', got %s", response.Metadata.Provider)
	}
	
	if response.Metadata.Details["total_words"] == "" {
		t.Error("Expected total_words in details")
	}
	
	if response.Metadata.Details["best_score"] == "" {
		t.Error("Expected best_score in details")
	}
}

func TestFallbackAdapter_CalculateLanguageScore(t *testing.T) {
	adapter := NewFallbackAdapter()
	
	tests := []struct {
		name     string
		words    []string
		patterns []string
		expected float32
	}{
		{
			name:     "No matches",
			words:    []string{"xyz", "abc", "def"},
			patterns: []string{"hello", "world"},
			expected: 0,
		},
		{
			name:     "Some matches",
			words:    []string{"hello", "world", "xyz"},
			patterns: []string{"hello", "world", "goodbye"},
			expected: 2,
		},
		{
			name:     "All matches",
			words:    []string{"hello", "world"},
			patterns: []string{"hello", "world"},
			expected: 2,
		},
		{
			name:     "Empty words",
			words:    []string{},
			patterns: []string{"hello", "world"},
			expected: 0,
		},
		{
			name:     "Empty patterns",
			words:    []string{"hello", "world"},
			patterns: []string{},
			expected: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := adapter.calculateLanguageScore(tt.words, tt.patterns)
			if score != tt.expected {
				t.Errorf("calculateLanguageScore() = %v, want %v", score, tt.expected)
			}
		})
	}
}

func TestFallbackAdapter_CalculateLanguageScore_WithPunctuation(t *testing.T) {
	adapter := NewFallbackAdapter()
	
	words := []string{"hello!", "world.", "test123"}
	patterns := []string{"hello", "world", "test"}
	
	score := adapter.calculateLanguageScore(words, patterns)
	
	// The score should be at least 2 (hello and world match)
	if score < 2 {
		t.Errorf("Expected score >= 2 with punctuation handling, got %v", score)
	}
}

func TestFallbackAdapter_ContextCancellation(t *testing.T) {
	adapter := NewFallbackAdapter()
	
	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	text := "Hello world"
	response, err := adapter.DetectLanguage(ctx, domain.Text(text))
	
	// The fallback adapter doesn't actually check context, so it should still work
	// This is expected behavior for this simple implementation
	if err != nil {
		t.Fatalf("Expected no error (context not checked), got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
}
