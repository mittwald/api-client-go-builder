package util

import "strings"

func UpperFirst(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}
