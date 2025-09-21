package domain

import "errors"

// Common domain errors for language detection
var (
	ErrEmptyText           = errors.New("text cannot be empty")
	ErrTextTooLong         = errors.New("text exceeds maximum length")
	ErrInvalidLanguageCode = errors.New("invalid language code")
	ErrLowConfidence       = errors.New("language detection confidence too low")
	ErrInvalidRequest      = errors.New("invalid request parameters")
	ErrInternalError       = errors.New("internal language detection error")
)
