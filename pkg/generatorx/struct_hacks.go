package generatorx

import (
	"github.com/mitchellh/go-wordwrap"
	"github.com/moznion/gowrtr/generator"
	"strings"
)

func AddFieldComment(str *generator.Struct, comment string) *generator.Struct {
	wrapped := wordwrap.WrapString(comment, 80)
	lines := strings.Split(wrapped, "\n")

	for _, l := range lines {
		str = str.AddField("//", l)
	}

	return str
}
