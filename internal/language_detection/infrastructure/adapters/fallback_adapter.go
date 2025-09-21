package adapters

import (
	"context"
	"fmt"
	"strings"

	"language-detection-service/internal/language_detection/domain"
)

// FallbackAdapter implements the LanguageDetector interface using pattern matching
type FallbackAdapter struct {
	patterns map[string][]string
}

// NewFallbackAdapter creates a new fallback adapter
func NewFallbackAdapter() *FallbackAdapter {
	return &FallbackAdapter{
		patterns: map[string][]string{
			"en-US": {
				"the", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by",
				"a", "an", "is", "are", "was", "were", "be", "been", "being", "have", "has", "had",
				"do", "does", "did", "will", "would", "could", "should", "may", "might", "can",
				"this", "that", "these", "those", "i", "you", "he", "she", "it", "we", "they",
				"welcome", "presentation", "today", "discussing", "technology", "trends", "talk",
				"artificial", "intelligence", "revolutionizing", "industries", "worldwide",
				"machine", "learning", "algorithms", "becoming", "sophisticated", "explore",
				"cloud", "computing", "solutions", "platforms", "offer", "scalable", "infrastructure",
				"finally", "examine", "cybersecurity", "measures", "protecting", "data", "crucial",
				"digital", "world", "thank", "attention",
			},
			"es-ES": {
				"hola", "bienvenidos", "presentación", "hoy", "vamos", "discutir", "sobre",
				"tendencias", "tecnológicas", "primero", "hablemos", "inteligencia", "artificial",
				"revolucionando", "muchas", "industrias", "mundo", "algoritmos", "aprendizaje",
				"automático", "volviendo", "sofisticados", "siguiente", "exploraremos", "soluciones",
				"computación", "nube", "plataformas", "ofrecen", "infraestructura", "escalable",
				"finalmente", "examinemos", "medidas", "ciberseguridad", "proteger", "datos",
				"crucial", "digital", "gracias", "atención", "el", "la", "los", "las", "de", "del",
				"que", "y", "en", "un", "una", "es", "son", "por", "para", "con", "sin",
			},
			"fr-FR": {
				"bonjour", "bienvenue", "présentation", "aujourd'hui", "allons", "discuter",
				"tendances", "technologiques", "intelligence", "artificielle", "révolutionne",
				"nombreuses", "industries", "algorithmes", "apprentissage", "automatique",
				"deviennent", "sophistiqués", "ensuite", "explorerons", "solutions", "informatique",
				"nuage", "plateformes", "offrent", "infrastructure", "évolutive", "enfin",
				"examinerons", "mesures", "cybersécurité", "protéger", "données", "crucial",
				"numérique", "merci", "attention", "le", "la", "les", "de", "du", "des", "que",
				"et", "en", "un", "une", "est", "sont", "pour", "par", "avec", "sans",
			},
			"de-DE": {
				"willkommen", "präsentation", "heute", "werden", "diskutieren", "neuesten",
				"technologietrends", "zuerst", "sprechen", "über", "künstliche", "intelligenz",
				"revolutioniert", "viele", "branchen", "weltweit", "algorithmen", "maschinelles",
				"lernen", "werden", "raffinierter", "als", "nächstes", "werden", "wir", "explorieren",
				"cloud-computing", "lösungen", "cloud-plattformen", "bieten", "skalierbare",
				"infrastruktur", "schließlich", "lassen", "uns", "cybersicherheitsmaßnahmen",
				"untersuchen", "datenschutz", "ist", "entscheidend", "in", "der", "heutigen",
				"digitalen", "welt", "danke", "für", "ihre", "aufmerksamkeit", "der", "die", "das",
				"und", "in", "auf", "zu", "für", "von", "mit", "ohne", "über", "unter",
			},
		},
	}
}

// DetectLanguage detects language using pattern matching
func (f *FallbackAdapter) DetectLanguage(
	ctx context.Context,
	text domain.Text,
) (*domain.LanguageDetectionResponse, error) {
	textStr := strings.ToLower(strings.TrimSpace(string(text)))

	if len(textStr) < 3 {
		return &domain.LanguageDetectionResponse{
			LanguageCode: domain.LanguageCode("unknown"),
			Confidence:   0,
			Metadata: domain.ProcessingMetadata{
				Provider: "fallback",
				Details: map[string]string{
					"reason": "text_too_short",
				},
			},
		}, nil
	}

	words := strings.Fields(textStr)
	if len(words) == 0 {
		return &domain.LanguageDetectionResponse{
			LanguageCode: domain.LanguageCode("unknown"),
			Confidence:   0,
			Metadata: domain.ProcessingMetadata{
				Provider: "fallback",
				Details: map[string]string{
					"reason": "no_words_found",
				},
			},
		}, nil
	}

	// Count words for each language
	languageScores := make(map[domain.LanguageCode]float32)
	totalWords := float32(len(words))

	for lang, patterns := range f.patterns {
		score := f.calculateLanguageScore(words, patterns)
		languageScores[domain.LanguageCode(lang)] = score / totalWords
	}

	// Find the language with highest score
	var bestLang domain.LanguageCode
	var bestScore float32

	for lang, score := range languageScores {
		if score > bestScore {
			bestLang = lang
			bestScore = score
		}
	}

	// Create alternatives
	var alternatives []domain.LanguageAlternative
	for lang, score := range languageScores {
		if lang != bestLang && score > 0.05 { // Only include alternatives with >5% confidence
			alternatives = append(alternatives, domain.LanguageAlternative{
				LanguageCode: lang,
				Confidence:   domain.Confidence(score),
			})
		}
	}

	// If no language meets threshold, return unknown
	if bestScore < 0.15 { // 15% threshold
		return &domain.LanguageDetectionResponse{
			LanguageCode: domain.LanguageCode("unknown"),
			Confidence:   domain.Confidence(bestScore),
			Alternatives: alternatives,
			Metadata: domain.ProcessingMetadata{
				Provider: "fallback",
				Details: map[string]string{
					"reason": "low_confidence",
					"best_score": fmt.Sprintf("%.3f", bestScore),
				},
			},
		}, nil
	}

	return &domain.LanguageDetectionResponse{
		LanguageCode: bestLang,
		Confidence:   domain.Confidence(bestScore),
		Alternatives: alternatives,
		Metadata: domain.ProcessingMetadata{
			Provider: "fallback",
			Details: map[string]string{
				"total_words": fmt.Sprintf("%d", len(words)),
				"best_score":  fmt.Sprintf("%.3f", bestScore),
			},
		},
	}, nil
}

// calculateLanguageScore calculates the score for a language based on word patterns
func (f *FallbackAdapter) calculateLanguageScore(words []string, patterns []string) float32 {
	score := float32(0)
	patternSet := make(map[string]bool)
	
	// Create pattern set for faster lookup
	for _, pattern := range patterns {
		patternSet[pattern] = true
	}

	// Count matching words
	for _, word := range words {
		// Clean word (remove punctuation)
		cleanWord := strings.ToLower(strings.TrimFunc(word, func(r rune) bool {
			return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
		}))

		if patternSet[cleanWord] {
			score++
		}
	}

	return score
}
