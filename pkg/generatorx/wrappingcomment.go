package generatorx

import (
	"fmt"
	"github.com/mitchellh/go-wordwrap"
)

type WrappingComment struct {
	comment string
}

func NewWrappingComment(comment string) *WrappingComment {
	return &WrappingComment{
		comment: comment,
	}
}

func NewWrappingCommentf(comment string, args ...any) *WrappingComment {
	return &WrappingComment{
		comment: fmt.Sprintf(comment, args...),
	}
}

func (c *WrappingComment) Writef(line string, args ...any) {
	c.comment += fmt.Sprintf(line, args...)
}

func (c *WrappingComment) Writeln(line string) {
	c.comment += "\n" + line
}

// Generate generates one line comment statement.
func (c *WrappingComment) Generate(indentLevel int) (string, error) {
	comment := wordwrap.WrapString(c.comment, 80)
	return NewMultilineComment(comment).Generate(indentLevel)
}
