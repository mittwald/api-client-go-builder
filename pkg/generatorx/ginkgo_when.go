package generatorx

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
	"strings"
)

type When struct {
	comment   string
	suiteFunc *generator.AnonymousFunc
}

func NewWhen(comment string, suiteFunc *generator.AnonymousFunc) *When {
	return &When{
		comment:   comment,
		suiteFunc: suiteFunc,
	}
}

// Generate generates one line comment statement.
func (c *When) Generate(indentLevel int) (string, error) {
	indent := generator.BuildIndent(indentLevel)

	suiteCode, err := c.suiteFunc.Generate(indentLevel + 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%sWhen(%#v, %s)", indent, c.comment, strings.Trim(suiteCode, " \t\n")), nil
}
