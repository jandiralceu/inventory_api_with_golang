package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "string with accents",
			input:    "Olá Café da Manhã",
			expected: "ola-cafe-da-manha",
		},
		{
			name:     "string with special characters",
			input:    "Product @ 2023! #Special",
			expected: "product-2023-special",
		},
		{
			name:     "multiple spaces and dashes",
			input:    "Multiple   Spaces --- and dashes",
			expected: "multiple-spaces-and-dashes",
		},
		{
			name:     "leading and trailing spaces",
			input:    "  Trim Me Please  ",
			expected: "trim-me-please",
		},
		{
			name:     "unicode characters",
			input:    "I ♥ Go 2024",
			expected: "i-go-2024",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!!!@@@###",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slugify(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
