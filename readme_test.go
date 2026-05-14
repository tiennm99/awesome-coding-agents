package main

import (
	"testing"
)

func TestSanitizeCell(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "plain text",
			input:    "hello world",
			expected: "hello world",
			desc:     "unchanged",
		},
		{
			name:     "pipe char",
			input:    "foo | bar",
			expected: "foo \\| bar",
			desc:     "pipe escaped",
		},
		{
			name:     "newline",
			input:    "line1\nline2",
			expected: "line1 line2",
			desc:     "newline becomes space",
		},
		{
			name:     "carriage return",
			input:    "line1\rline2",
			expected: "line1 line2",
			desc:     "carriage return becomes space",
		},
		{
			name:     "pipe and newline",
			input:    "foo | bar\nbaz",
			expected: "foo \\| bar baz",
			desc:     "both handled correctly",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			desc:     "empty string unchanged",
		},
		{
			name:     "multiple pipes",
			input:    "a | b | c",
			expected: "a \\| b \\| c",
			desc:     "multiple pipes escaped",
		},
		{
			name:     "multiple newlines",
			input:    "line1\n\nline2",
			expected: "line1  line2",
			desc:     "multiple newlines become spaces",
		},
		{
			name:     "whitespace trimming",
			input:    "  text  ",
			expected: "text",
			desc:     "leading/trailing spaces trimmed",
		},
		{
			name:     "complex: whitespace, pipe, newline",
			input:    "  foo | bar\nbaz  ",
			expected: "foo \\| bar baz",
			desc:     "all transformations applied",
		},
		{
			name:     "only newlines and pipes",
			input:    "|\n|",
			expected: "\\| \\|",
			desc:     "pipes and newlines only",
		},
		{
			name:     "mixed line endings",
			input:    "line1\nline2\rline3",
			expected: "line1 line2 line3",
			desc:     "both \\n and \\r converted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeCell(tt.input)
			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.desc, tt.expected, result)
			}
		})
	}
}
