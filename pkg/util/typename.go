package util

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func ConvertToTypename(input string) string {
	caser := cases.Title(language.English)
	upper := caser.String(input)

	cleaned := upper
	cleaned = strings.Replace(cleaned, "-", "", -1)
	cleaned = strings.Replace(cleaned, "_", "", -1)
	cleaned = strings.Replace(cleaned, " ", "", -1)

	return cleaned
}
