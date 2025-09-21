package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"language-detection-service/internal/language_detection/domain"
)

// Config holds all configuration for the language detection service
type Config struct {
	// Server configuration
	ServerAddress string
	ServerPort    int

	// AWS configuration
	AWSRegion        string
	UseAWSComprehend bool

	// Service configuration
	MaxTextLength          int
	MinConfidenceThreshold float32
	ServiceVersion         string
	ModelVersion           string

	// Supported languages
	SupportedLanguages []domain.LanguageCode

	// Timeouts
	ShutdownTimeoutSeconds int
}

// ConfigProvider implements the domain.ConfigProvider interface
type ConfigProvider struct {
	config *Config
}

// NewConfigProvider creates a new configuration provider
func NewConfigProvider() *ConfigProvider {
	config := &Config{
		ServerAddress:          getEnv("SERVER_ADDRESS", "0.0.0.0"),
		ServerPort:             getEnvInt("SERVER_PORT", 6011),
		AWSRegion:              getEnv("AWS_REGION", "us-east-1"),
		UseAWSComprehend:       getEnvBool("USE_AWS_COMPREHEND", true),
		MaxTextLength:          getEnvInt("MAX_TEXT_LENGTH", 5000),
		MinConfidenceThreshold: getEnvFloat32("MIN_CONFIDENCE_THRESHOLD", 0.1),
		ServiceVersion:         getEnv("SERVICE_VERSION", "1.0.0"),
		ModelVersion:           getEnv("MODEL_VERSION", "1.0.0"),
		ShutdownTimeoutSeconds: getEnvInt("SHUTDOWN_TIMEOUT_SECONDS", 30),
	}

	// Parse supported languages
	config.SupportedLanguages = parseSupportedLanguages(getEnv("SUPPORTED_LANGUAGES", ""))

	return &ConfigProvider{config: config}
}

// GetConfig returns the underlying configuration
func (cp *ConfigProvider) GetConfig() *Config {
	return cp.config
}

// Implement domain.ConfigProvider interface
func (cp *ConfigProvider) GetMaxTextLength() int {
	return cp.config.MaxTextLength
}

func (cp *ConfigProvider) GetMinConfidenceThreshold() float32 {
	return cp.config.MinConfidenceThreshold
}

func (cp *ConfigProvider) GetSupportedLanguages() []domain.LanguageCode {
	return cp.config.SupportedLanguages
}

func (cp *ConfigProvider) GetServiceVersion() string {
	return cp.config.ServiceVersion
}

func (cp *ConfigProvider) GetModelVersion() string {
	return cp.config.ModelVersion
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat32(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(floatValue)
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func parseSupportedLanguages(languagesStr string) []domain.LanguageCode {
	if languagesStr == "" {
		// Default supported languages
		return []domain.LanguageCode{
			"en-US", "es-ES", "fr-FR", "de-DE", "it-IT", "pt-PT", "ru-RU",
			"ja-JP", "ko-KR", "zh-CN", "ar-SA", "hi-IN", "unknown",
		}
	}

	languages := strings.Split(languagesStr, ",")
	var result []domain.LanguageCode
	for _, lang := range languages {
		lang = strings.TrimSpace(lang)
		if lang != "" {
			result = append(result, domain.LanguageCode(lang))
		}
	}

	if len(result) == 0 {
		// Fallback to default if parsing failed
		return []domain.LanguageCode{"en-US", "unknown"}
	}

	return result
}

// ValidateConfig validates the configuration
func (cp *ConfigProvider) ValidateConfig() error {
	config := cp.config

	// Validate server configuration
	if config.ServerPort <= 0 || config.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", config.ServerPort)
	}

	// Validate AWS configuration if using AWS Comprehend
	if config.UseAWSComprehend {
		if config.AWSRegion == "" {
			return fmt.Errorf("AWS region is required when using AWS Comprehend")
		}
	}

	// Validate text length
	if config.MaxTextLength <= 0 {
		return fmt.Errorf("max text length must be positive")
	}

	// Validate confidence threshold
	if config.MinConfidenceThreshold < 0 || config.MinConfidenceThreshold > 1 {
		return fmt.Errorf("confidence threshold must be between 0 and 1")
	}

	// Validate supported languages
	if len(config.SupportedLanguages) == 0 {
		return fmt.Errorf("at least one supported language must be configured")
	}

	// Validate timeouts
	if config.ShutdownTimeoutSeconds <= 0 {
		return fmt.Errorf("shutdown timeout must be positive")
	}

	return nil
}
