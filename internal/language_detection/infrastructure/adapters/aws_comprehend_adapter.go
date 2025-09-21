package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"

	"language-detection-service/internal/language_detection/domain"
)

// AWSComprehendAdapter implements the LanguageDetector interface using AWS Comprehend
type AWSComprehendAdapter struct {
	client     *comprehend.Comprehend
	region     string
	maxRetries int
}

// NewAWSComprehendAdapter creates a new AWS Comprehend adapter
func NewAWSComprehendAdapter(region string, maxRetries int) (*AWSComprehendAdapter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := comprehend.New(sess)

	return &AWSComprehendAdapter{
		client:     client,
		region:     region,
		maxRetries: maxRetries,
	}, nil
}

// DetectLanguage detects language using AWS Comprehend
func (a *AWSComprehendAdapter) DetectLanguage(
	ctx context.Context,
	text domain.Text,
) (*domain.LanguageDetectionResponse, error) {
	textStr := string(text)

	// Truncate text if too long (Comprehend has a 5000 character limit per document)
	if len(textStr) > 5000 {
		textStr = textStr[:5000]
	}

	// If text is empty or too short, return unknown
	if len(strings.TrimSpace(textStr)) < 3 {
		return &domain.LanguageDetectionResponse{
			LanguageCode: domain.LanguageCode("unknown"),
			Confidence:   0,
			Metadata: domain.ProcessingMetadata{
				Provider: "aws-comprehend",
				Details: map[string]string{
					"reason": "text_too_short",
				},
			},
		}, nil
	}

	// Call AWS Comprehend
	input := &comprehend.DetectDominantLanguageInput{
		Text: aws.String(textStr),
	}

	result, err := a.client.DetectDominantLanguageWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("AWS Comprehend error: %w", err)
	}

	// Get the most confident language
	if len(result.Languages) == 0 {
		return &domain.LanguageDetectionResponse{
			LanguageCode: domain.LanguageCode("unknown"),
			Confidence:   0,
			Metadata: domain.ProcessingMetadata{
				Provider: "aws-comprehend",
				Details: map[string]string{
					"reason": "no_languages_detected",
				},
			},
		}, nil
	}

	// Find the language with highest confidence
	var dominantLang *comprehend.DominantLanguage
	for _, lang := range result.Languages {
		if dominantLang == nil || *lang.Score > *dominantLang.Score {
			dominantLang = lang
		}
	}

	if dominantLang == nil || dominantLang.LanguageCode == nil {
		return &domain.LanguageDetectionResponse{
			LanguageCode: domain.LanguageCode("unknown"),
			Confidence:   0,
			Metadata: domain.ProcessingMetadata{
				Provider: "aws-comprehend",
				Details: map[string]string{
					"reason": "invalid_response",
				},
			},
		}, nil
	}

	// Convert AWS language code to our format
	langCode := a.convertLanguageCode(*dominantLang.LanguageCode)
	confidence := domain.Confidence(*dominantLang.Score)

	// Create alternatives from other detected languages
	var alternatives []domain.LanguageAlternative
	for _, lang := range result.Languages {
		if lang.LanguageCode != nil && *lang.LanguageCode != *dominantLang.LanguageCode {
			alternatives = append(alternatives, domain.LanguageAlternative{
				LanguageCode: a.convertLanguageCode(*lang.LanguageCode),
				Confidence:   domain.Confidence(*lang.Score),
			})
		}
	}

	return &domain.LanguageDetectionResponse{
		LanguageCode: langCode,
		Confidence:   confidence,
		Alternatives: alternatives,
		Metadata: domain.ProcessingMetadata{
			Provider: "aws-comprehend",
			Details: map[string]string{
				"region":       a.region,
				"total_langs":  fmt.Sprintf("%d", len(result.Languages)),
				"aws_lang_code": *dominantLang.LanguageCode,
			},
		},
	}, nil
}

// convertLanguageCode converts AWS language codes to our standard format
func (a *AWSComprehendAdapter) convertLanguageCode(awsCode string) domain.LanguageCode {
	langCode := strings.ToLower(awsCode)
	
	// Map AWS language codes to our expected format
	switch langCode {
	case "en":
		return "en-US"
	case "es":
		return "es-ES"
	case "fr":
		return "fr-FR"
	case "de":
		return "de-DE"
	case "it":
		return "it-IT"
	case "pt":
		return "pt-PT"
	case "ru":
		return "ru-RU"
	case "ja":
		return "ja-JP"
	case "ko":
		return "ko-KR"
	case "zh", "zh-tw", "zh-cn":
		return "zh-CN"
	case "ar":
		return "ar-SA"
	case "hi":
		return "hi-IN"
	default:
		return domain.LanguageCode(awsCode)
	}
}
