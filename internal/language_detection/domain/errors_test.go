package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "EmptyText error",
			err:      ErrEmptyText,
			expected: "text cannot be empty",
		},
		{
			name:     "TextTooLong error",
			err:      ErrTextTooLong,
			expected: "text exceeds maximum length",
		},
		{
			name:     "InvalidLanguageCode error",
			err:      ErrInvalidLanguageCode,
			expected: "invalid language code",
		},
		{
			name:     "LowConfidence error",
			err:      ErrLowConfidence,
			expected: "language detection confidence too low",
		},
		{
			name:     "InvalidRequest error",
			err:      ErrInvalidRequest,
			expected: "invalid request parameters",
		},
		{
			name:     "InternalError error",
			err:      ErrInternalError,
			expected: "internal language detection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Error message = %v, want %v", tt.err.Error(), tt.expected)
			}
		})
	}
}

func TestErrorComparisons(t *testing.T) {
	tests := []struct {
		name     string
		err1     error
		err2     error
		expected bool
	}{
		{
			name:     "Same error type",
			err1:     ErrEmptyText,
			err2:     ErrEmptyText,
			expected: true,
		},
		{
			name:     "Different error types",
			err1:     ErrEmptyText,
			err2:     ErrTextTooLong,
			expected: false,
		},
		{
			name:     "Error vs nil",
			err1:     ErrEmptyText,
			err2:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err1, tt.err2)
			if result != tt.expected {
				t.Errorf("errors.Is(%v, %v) = %v, want %v", tt.err1, tt.err2, result, tt.expected)
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Wrap ErrEmptyText",
			err:      ErrEmptyText,
			expected: true,
		},
		{
			name:     "Wrap ErrTextTooLong",
			err:      ErrTextTooLong,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedErr := fmt.Errorf("wrapper: %w", tt.err)
			if !errors.Is(wrappedErr, tt.err) {
				t.Errorf("Wrapped error should contain original error")
			}
		})
	}
}
