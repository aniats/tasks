package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		escapeEnabled bool
		expected      string
		shouldError   bool
	}{
		{
			name:          "basic unpacking",
			input:         "a4bc2d5e",
			escapeEnabled: false,
			expected:      "aaaabccddddde",
			shouldError:   false,
		},
		{
			name:          "no numbers",
			input:         "abcd",
			escapeEnabled: false,
			expected:      "abcd",
			shouldError:   false,
		},
		{
			name:          "starts with number",
			input:         "3abc",
			escapeEnabled: false,
			expected:      "",
			shouldError:   true,
		},
		{
			name:          "only numbers",
			input:         "45",
			escapeEnabled: false,
			expected:      "",
			shouldError:   true,
		},
		{
			name:          "consecutive numbers",
			input:         "aaa10b",
			escapeEnabled: false,
			expected:      "",
			shouldError:   true,
		},
		{
			name:          "zero removes character",
			input:         "aaa0b",
			escapeEnabled: false,
			expected:      "aab",
			shouldError:   false,
		},
		{
			name:          "empty string",
			input:         "",
			escapeEnabled: false,
			expected:      "",
			shouldError:   false,
		},
		{
			name:          "newline character",
			input:         "d\n5abc",
			escapeEnabled: false,
			expected:      "d\n\n\n\n\nabc",
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := unpackString(tt.input, tt.escapeEnabled)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got result: %q", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				} else if result != tt.expected {
					t.Errorf("For input %q, expected %q, but got %q", tt.input, tt.expected, result)
				}
			}
		})
	}
}

func TestUnpackStringWithEscape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expected      string
		escapeEnabled bool
		shouldError   bool
	}{
		{
			name:          `qwe\4\5 => qwe45`,
			input:         `qwe\4\5`,
			expected:      `qwe45`,
			shouldError:   false,
			escapeEnabled: true,
		},
		{
			name:          `qwe\45 => qwe44444`,
			input:         `qwe\45`,
			expected:      `qwe44444`,
			shouldError:   false,
			escapeEnabled: true,
		},
		{
			name:          `qwe\\5 => qwe\\\\\\\\\\`,
			input:         `qwe\\5`,
			expected:      `qwe\\\\\`,
			shouldError:   false,
			escapeEnabled: true,
		},
		{
			name:          `qw\ne => error`,
			input:         `qw\ne`,
			expected:      ``,
			shouldError:   true,
			escapeEnabled: true,
		},
		{
			name:          `backslash at end`,
			input:         `abc\\`,
			expected:      ``,
			shouldError:   true,
			escapeEnabled: true,
		},
		{
			name:          `escape zero`,
			input:         `ab\\0c`,
			expected:      `abc`,
			shouldError:   false,
			escapeEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := unpackString(tt.input, tt.escapeEnabled)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got result: %q", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				} else if result != tt.expected {
					t.Errorf("For input %q, expected %q, but got %q", tt.input, tt.expected, result)
				}
			}
		})
	}
}

func TestPackString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		expected    string
		shouldError bool
	}{
		{
			name:        "basic packing",
			input:       "aaaabccddddde",
			expected:    "a4bc2d5e",
			shouldError: false,
		},
		{
			name:        "no repetitions",
			input:       "abcd",
			expected:    "abcd",
			shouldError: false,
		},
		{
			name:        "single character",
			input:       "a",
			expected:    "a",
			shouldError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			shouldError: false,
		},
		{
			name:        "all same characters",
			input:       "aaaa",
			expected:    "a4",
			shouldError: false,
		},
		{
			name:        "newline characters",
			input:       "d\n\n\n\n\nabc",
			expected:    "d\n5abc",
			shouldError: false,
		},
		{
			name:        "mixed repetitions",
			input:       "aabbbbcccccccc",
			expected:    "a2b4c8",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := packString(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got result: %q", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				} else if result != tt.expected {
					t.Errorf("For input %q, expected %q, but got %q", tt.input, tt.expected, result)
				}
			}
		})
	}
}
