package config

import (
	"os"
	"testing"

	"language-detection-service/internal/language_detection/domain"
)

func TestNewConfigProvider(t *testing.T) {
	provider := NewConfigProvider()
	
	if provider == nil {
		t.Fatal("Expected provider to be created, got nil")
	}
	
	if provider.config == nil {
		t.Fatal("Expected config to be initialized, got nil")
	}
}

func TestConfigProvider_DefaultValues(t *testing.T) {
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	// Test default values
	if config.ServerAddress != "0.0.0.0" {
		t.Errorf("Expected ServerAddress '0.0.0.0', got %s", config.ServerAddress)
	}
	
	if config.ServerPort != 6011 {
		t.Errorf("Expected ServerPort 6011, got %d", config.ServerPort)
	}
	
	if config.AWSRegion != "us-east-1" {
		t.Errorf("Expected AWSRegion 'us-east-1', got %s", config.AWSRegion)
	}
	
	if config.UseAWSComprehend != true {
		t.Errorf("Expected UseAWSComprehend true, got %v", config.UseAWSComprehend)
	}
	
	if config.MaxTextLength != 5000 {
		t.Errorf("Expected MaxTextLength 5000, got %d", config.MaxTextLength)
	}
	
	if config.MinConfidenceThreshold != 0.1 {
		t.Errorf("Expected MinConfidenceThreshold 0.1, got %f", config.MinConfidenceThreshold)
	}
	
	if config.ServiceVersion != "1.0.0" {
		t.Errorf("Expected ServiceVersion '1.0.0', got %s", config.ServiceVersion)
	}
	
	if config.ModelVersion != "1.0.0" {
		t.Errorf("Expected ModelVersion '1.0.0', got %s", config.ModelVersion)
	}
	
	if config.ShutdownTimeoutSeconds != 30 {
		t.Errorf("Expected ShutdownTimeoutSeconds 30, got %d", config.ShutdownTimeoutSeconds)
	}
}

func TestConfigProvider_InterfaceMethods(t *testing.T) {
	provider := NewConfigProvider()
	
	// Test interface methods
	maxLength := provider.GetMaxTextLength()
	if maxLength != 5000 {
		t.Errorf("GetMaxTextLength() = %d, want 5000", maxLength)
	}
	
	minConfidence := provider.GetMinConfidenceThreshold()
	if minConfidence != 0.1 {
		t.Errorf("GetMinConfidenceThreshold() = %f, want 0.1", minConfidence)
	}
	
	supportedLangs := provider.GetSupportedLanguages()
	if len(supportedLangs) == 0 {
		t.Error("Expected supported languages, got empty slice")
	}
	
	// Check that default languages are present
	expectedLangs := []domain.LanguageCode{"en-US", "es-ES", "fr-FR", "de-DE", "it-IT", "pt-PT", "ru-RU", "ja-JP", "ko-KR", "zh-CN", "ar-SA", "hi-IN", "unknown"}
	langMap := make(map[domain.LanguageCode]bool)
	for _, lang := range supportedLangs {
		langMap[lang] = true
	}
	
	for _, expectedLang := range expectedLangs {
		if !langMap[expectedLang] {
			t.Errorf("Expected supported language %s not found", expectedLang)
		}
	}
	
	serviceVersion := provider.GetServiceVersion()
	if serviceVersion != "1.0.0" {
		t.Errorf("GetServiceVersion() = %s, want '1.0.0'", serviceVersion)
	}
	
	modelVersion := provider.GetModelVersion()
	if modelVersion != "1.0.0" {
		t.Errorf("GetModelVersion() = %s, want '1.0.0'", modelVersion)
	}
}

func TestConfigProvider_EnvironmentVariables(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"SERVER_ADDRESS", "SERVER_PORT", "AWS_REGION", "USE_AWS_COMPREHEND",
		"MAX_TEXT_LENGTH", "MIN_CONFIDENCE_THRESHOLD", "SERVICE_VERSION",
		"MODEL_VERSION", "SHUTDOWN_TIMEOUT_SECONDS", "SUPPORTED_LANGUAGES",
	}
	
	for _, envVar := range envVars {
		originalEnv[envVar] = os.Getenv(envVar)
	}
	
	// Clean up after test
	defer func() {
		for envVar, value := range originalEnv {
			if value == "" {
				os.Unsetenv(envVar)
			} else {
				os.Setenv(envVar, value)
			}
		}
	}()
	
	// Set test environment variables
	os.Setenv("SERVER_ADDRESS", "127.0.0.1")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("USE_AWS_COMPREHEND", "false")
	os.Setenv("MAX_TEXT_LENGTH", "10000")
	os.Setenv("MIN_CONFIDENCE_THRESHOLD", "0.5")
	os.Setenv("SERVICE_VERSION", "2.0.0")
	os.Setenv("MODEL_VERSION", "2.0.0")
	os.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "60")
	os.Setenv("SUPPORTED_LANGUAGES", "en-US,es-ES,fr-FR")
	
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	// Test environment variable values
	if config.ServerAddress != "127.0.0.1" {
		t.Errorf("Expected ServerAddress '127.0.0.1', got %s", config.ServerAddress)
	}
	
	if config.ServerPort != 8080 {
		t.Errorf("Expected ServerPort 8080, got %d", config.ServerPort)
	}
	
	if config.AWSRegion != "us-west-2" {
		t.Errorf("Expected AWSRegion 'us-west-2', got %s", config.AWSRegion)
	}
	
	if config.UseAWSComprehend != false {
		t.Errorf("Expected UseAWSComprehend false, got %v", config.UseAWSComprehend)
	}
	
	if config.MaxTextLength != 10000 {
		t.Errorf("Expected MaxTextLength 10000, got %d", config.MaxTextLength)
	}
	
	if config.MinConfidenceThreshold != 0.5 {
		t.Errorf("Expected MinConfidenceThreshold 0.5, got %f", config.MinConfidenceThreshold)
	}
	
	if config.ServiceVersion != "2.0.0" {
		t.Errorf("Expected ServiceVersion '2.0.0', got %s", config.ServiceVersion)
	}
	
	if config.ModelVersion != "2.0.0" {
		t.Errorf("Expected ModelVersion '2.0.0', got %s", config.ModelVersion)
	}
	
	if config.ShutdownTimeoutSeconds != 60 {
		t.Errorf("Expected ShutdownTimeoutSeconds 60, got %d", config.ShutdownTimeoutSeconds)
	}
	
	// Test supported languages
	supportedLangs := provider.GetSupportedLanguages()
	expectedLangs := []domain.LanguageCode{"en-US", "es-ES", "fr-FR"}
	if len(supportedLangs) != len(expectedLangs) {
		t.Errorf("Expected %d supported languages, got %d", len(expectedLangs), len(supportedLangs))
	}
	
	for i, expected := range expectedLangs {
		if supportedLangs[i] != expected {
			t.Errorf("Expected supported language %s at index %d, got %s", expected, i, supportedLangs[i])
		}
	}
}

func TestConfigProvider_InvalidEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{"SERVER_PORT", "MAX_TEXT_LENGTH", "MIN_CONFIDENCE_THRESHOLD", "USE_AWS_COMPREHEND"}
	
	for _, envVar := range envVars {
		originalEnv[envVar] = os.Getenv(envVar)
	}
	
	// Clean up after test
	defer func() {
		for envVar, value := range originalEnv {
			if value == "" {
				os.Unsetenv(envVar)
			} else {
				os.Setenv(envVar, value)
			}
		}
	}()
	
	// Set invalid environment variables
	os.Setenv("SERVER_PORT", "invalid")
	os.Setenv("MAX_TEXT_LENGTH", "not_a_number")
	os.Setenv("MIN_CONFIDENCE_THRESHOLD", "invalid")
	os.Setenv("USE_AWS_COMPREHEND", "maybe")
	
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	// Should fall back to default values
	if config.ServerPort != 6011 {
		t.Errorf("Expected default ServerPort 6011, got %d", config.ServerPort)
	}
	
	if config.MaxTextLength != 5000 {
		t.Errorf("Expected default MaxTextLength 5000, got %d", config.MaxTextLength)
	}
	
	if config.MinConfidenceThreshold != 0.1 {
		t.Errorf("Expected default MinConfidenceThreshold 0.1, got %f", config.MinConfidenceThreshold)
	}
	
	if config.UseAWSComprehend != true {
		t.Errorf("Expected default UseAWSComprehend true, got %v", config.UseAWSComprehend)
	}
}

func TestParseSupportedLanguages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []domain.LanguageCode
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []domain.LanguageCode{"en-US", "es-ES", "fr-FR", "de-DE", "it-IT", "pt-PT", "ru-RU", "ja-JP", "ko-KR", "zh-CN", "ar-SA", "hi-IN", "unknown"},
		},
		{
			name:     "Single language",
			input:    "en-US",
			expected: []domain.LanguageCode{"en-US"},
		},
		{
			name:     "Multiple languages",
			input:    "en-US,es-ES,fr-FR",
			expected: []domain.LanguageCode{"en-US", "es-ES", "fr-FR"},
		},
		{
			name:     "Languages with spaces",
			input:    "en-US, es-ES , fr-FR ",
			expected: []domain.LanguageCode{"en-US", "es-ES", "fr-FR"},
		},
		{
			name:     "Empty entries",
			input:    "en-US,,es-ES,",
			expected: []domain.LanguageCode{"en-US", "es-ES"},
		},
		{
			name:     "Invalid parsing",
			input:    ",,,",
			expected: []domain.LanguageCode{"en-US", "unknown"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSupportedLanguages(tt.input)
			
			if len(result) != len(tt.expected) {
				t.Errorf("parseSupportedLanguages() returned %d languages, want %d", len(result), len(tt.expected))
			}
			
			for i, expected := range tt.expected {
				if i < len(result) && result[i] != expected {
					t.Errorf("parseSupportedLanguages()[%d] = %s, want %s", i, result[i], expected)
				}
			}
		})
	}
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	provider := NewConfigProvider()
	
	err := provider.ValidateConfig()
	if err != nil {
		t.Errorf("ValidateConfig() error = %v, want nil", err)
	}
}

func TestValidateConfig_InvalidServerPort(t *testing.T) {
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	// Test invalid port numbers
	tests := []struct {
		name string
		port int
	}{
		{"Negative port", -1},
		{"Zero port", 0},
		{"Port too high", 65536},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalPort := config.ServerPort
			config.ServerPort = tt.port
			
			err := provider.ValidateConfig()
			if err == nil {
				t.Errorf("ValidateConfig() expected error for port %d, got nil", tt.port)
			}
			
			// Restore original port
			config.ServerPort = originalPort
		})
	}
}

func TestValidateConfig_AWSComprehendWithoutRegion(t *testing.T) {
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	config.UseAWSComprehend = true
	config.AWSRegion = ""
	
	err := provider.ValidateConfig()
	if err == nil {
		t.Error("ValidateConfig() expected error for AWS Comprehend without region, got nil")
	}
}

func TestValidateConfig_InvalidMaxTextLength(t *testing.T) {
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	config.MaxTextLength = 0
	
	err := provider.ValidateConfig()
	if err == nil {
		t.Error("ValidateConfig() expected error for zero max text length, got nil")
	}
}

func TestValidateConfig_InvalidConfidenceThreshold(t *testing.T) {
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	tests := []struct {
		name      string
		threshold float32
	}{
		{"Negative threshold", -0.1},
		{"Threshold too high", 1.1},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalThreshold := config.MinConfidenceThreshold
			config.MinConfidenceThreshold = tt.threshold
			
			err := provider.ValidateConfig()
			if err == nil {
				t.Errorf("ValidateConfig() expected error for threshold %f, got nil", tt.threshold)
			}
			
			// Restore original threshold
			config.MinConfidenceThreshold = originalThreshold
		})
	}
}

func TestValidateConfig_EmptySupportedLanguages(t *testing.T) {
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	config.SupportedLanguages = []domain.LanguageCode{}
	
	err := provider.ValidateConfig()
	if err == nil {
		t.Error("ValidateConfig() expected error for empty supported languages, got nil")
	}
}

func TestValidateConfig_InvalidShutdownTimeout(t *testing.T) {
	provider := NewConfigProvider()
	config := provider.GetConfig()
	
	config.ShutdownTimeoutSeconds = 0
	
	err := provider.ValidateConfig()
	if err == nil {
		t.Error("ValidateConfig() expected error for zero shutdown timeout, got nil")
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test getEnv with default
	result := getEnv("NONEXISTENT_VAR", "default")
	if result != "default" {
		t.Errorf("getEnv() = %s, want 'default'", result)
	}
	
	// Test getEnvInt with default
	resultInt := getEnvInt("NONEXISTENT_INT", 42)
	if resultInt != 42 {
		t.Errorf("getEnvInt() = %d, want 42", resultInt)
	}
	
	// Test getEnvFloat32 with default
	resultFloat := getEnvFloat32("NONEXISTENT_FLOAT", 3.14)
	if resultFloat != 3.14 {
		t.Errorf("getEnvFloat32() = %f, want 3.14", resultFloat)
	}
	
	// Test getEnvBool with default
	resultBool := getEnvBool("NONEXISTENT_BOOL", true)
	if resultBool != true {
		t.Errorf("getEnvBool() = %v, want true", resultBool)
	}
}
