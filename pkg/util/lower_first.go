package util

import "strings"

func LowerFirst(input string) string {
	return strings.ToLower(input[0:1]) + input[1:]
}
