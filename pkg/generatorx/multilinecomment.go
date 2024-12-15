package generatorx

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
	"strings"
)

type MultilineComment struct {
	comment string
}

func NewMultilineComment(comment string) *MultilineComment {
	return &MultilineComment{
		comment: comment,
	}
}

// Generate generates one line comment statement.
func (c *MultilineComment) Generate(indentLevel int) (string, error) {
	indent := generator.BuildIndent(indentLevel)
	lines := strings.Split(c.comment, "\n")
	commentLines := make([]string, len(lines))

	for i, line := range lines {
		commentLines[i] = fmt.Sprintf("%s//%s\n", indent, line)
	}

	return strings.Join(commentLines, ""), nil
}
