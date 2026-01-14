package util

import (
	"fmt"
	"testing"
)

func TestConvertToTypename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"order-id", "OrderID"},
		{"customer_id", "CustomerID"},
		{"article.id", "ArticleID"},
		{"use_free_trial", "UseFreeTrial"},
		{"sftp-user", "SFTPUser"},
		{"ssh_user", "SSHUser"},
		{"api_endpoint", "APIEndpoint"},
		{"AI hosting", "AIHosting"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%q should convert to %q", test.input, test.expected), func(t *testing.T) {
			result := ConvertToTypename(test.input)
			if result != test.expected {
				t.Errorf("ConvertToTypename(%q) = %q; want %q", test.input, result, test.expected)
			}
		})
	}
}
