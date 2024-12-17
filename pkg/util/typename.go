package util

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func ConvertToTypename(input string) string {
	caser := cases.Title(language.English, cases.NoLower)
	upper := caser.String(input)

	cleaned := upper
	cleaned = strings.Replace(cleaned, "-", "", -1)
	cleaned = strings.Replace(cleaned, ";", "", -1)
	cleaned = strings.Replace(cleaned, "_", "", -1)
	cleaned = strings.Replace(cleaned, " ", "", -1)
	cleaned = strings.Replace(cleaned, "/", "", -1)

	return cleaned
}

func ConvertToPackagename(input string) string {
	lower := strings.ToLower(input)

	cleaned := lower
	cleaned = strings.Replace(cleaned, "-", "_", -1)
	cleaned = strings.Replace(cleaned, " ", "_", -1)
	cleaned = strings.Replace(cleaned, "/", "_", -1)

	return cleaned
}
