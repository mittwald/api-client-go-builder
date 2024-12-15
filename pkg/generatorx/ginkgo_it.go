package generatorx

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
	"strings"
)

type It struct {
	comment   string
	suiteFunc *generator.AnonymousFunc
}

func NewIt(comment string, suiteFunc *generator.AnonymousFunc) *It {
	return &It{
		comment:   comment,
		suiteFunc: suiteFunc,
	}
}

// Generate generates one line comment statement.
func (c *It) Generate(indentLevel int) (string, error) {
	indent := generator.BuildIndent(indentLevel)

	suiteCode, err := c.suiteFunc.Generate(indentLevel + 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%sIt(%#v, %s)\n", indent, c.comment, strings.Trim(suiteCode, " \t\n")), nil
}
