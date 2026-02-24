package llm

import (
	"testing"
)

func TestSanitizeJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Plain JSON",
			input:    `{"test": "value"}`,
			expected: `{"test": "value"}`,
		},
		{
			name:     "Markdown JSON block",
			input:    "```json\n{\"test\": \"value\"}\n```",
			expected: `{"test": "value"}`,
		},
		{
			name:     "Markdown code block without tag",
			input:    "```\n{\"test\": \"value\"}\n```",
			expected: `{"test": "value"}`,
		},
		{
			name:     "JSON with surrounding whitespace",
			input:    "   \n{\"test\": \"value\"}\t ",
			expected: `{"test": "value"}`,
		},
		{
			name:     "JSON within text",
			input:    "Here is the result: ```json\n{\"test\": \"value\"}\n```",
			expected: `{"test": "value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeJSON(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeJSON() = %q, want %q", got, tt.expected)
			}
		})
	}
}
