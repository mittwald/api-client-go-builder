package util

import (
	"strings"
	"unicode"
)

// ConvertToFieldname converts a JSON property name into an exported Go struct
// field name. Unlike ConvertToTypename, names that are already valid
// identifiers are only upper-cased on their first letter (no initialism
// handling), so that existing field names remain stable. Characters that are
// not valid in Go identifiers act as word separators.
func ConvertToFieldname(input string) string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
	})

	for i, p := range parts {
		parts[i] = UpperFirst(p)
	}

	return strings.Join(parts, "")
}
