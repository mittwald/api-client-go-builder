package util

import (
	"fmt"
	"testing"
)

func TestConvertToFieldname(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"projectId", "ProjectId"},
		{"description", "Description"},
		{"use_free_trial", "Use_free_trial"},
		{"x-description", "XDescription"},
		{"x-update-schedule", "XUpdateSchedule"},
		{"foo.bar", "FooBar"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%q should convert to %q", test.input, test.expected), func(t *testing.T) {
			result := ConvertToFieldname(test.input)
			if result != test.expected {
				t.Errorf("ConvertToFieldname(%q) = %q; want %q", test.input, result, test.expected)
			}
		})
	}
}
