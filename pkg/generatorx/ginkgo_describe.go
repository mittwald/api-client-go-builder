package generatorx

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
	"strings"
)

type Describe struct {
	comment   string
	suiteFunc *generator.AnonymousFunc
}

func NewDescribe(comment string, suiteFunc *generator.AnonymousFunc) *Describe {
	return &Describe{
		comment:   comment,
		suiteFunc: suiteFunc,
	}
}

// Generate generates one line comment statement.
func (c *Describe) Generate(indentLevel int) (string, error) {
	indent := generator.BuildIndent(indentLevel)

	suiteCode, err := c.suiteFunc.Generate(indentLevel + 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%svar _ = Describe(%#v, %s)", indent, c.comment, strings.Trim(suiteCode, " \t\n")), nil
}
