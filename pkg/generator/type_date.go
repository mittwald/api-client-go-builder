package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/moznion/gowrtr/generator"
	"strings"
	"time"
)

var _ SchemaType = &DateType{}

type DateType struct {
	BaseType
}

func (o *DateType) IsLightweight() bool {
	return true
}

func (o *DateType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	stmts := make([]generator.Statement, 0)

	if d := o.schema.Schema().Description; d != "" {
		stmts = append(stmts, generatorx.NewMultilineComment(d))
	}

	stmts = append(stmts, generator.NewRawStatementf("type %s time.Time", o.Names.StructName))
	return stmts
}

func (o *DateType) EmitReference(ctx *GeneratorContext) string {
	return "time.Time"
}

func (o *DateType) BuildExample(*GeneratorContext, int, int) any {
	if example := o.schema.Schema().Example; example != nil {
		return example.Value
	}

	if examples := o.schema.Schema().Examples; len(examples) > 0 {
		return examples[0].Value
	}

	example, _ := time.Parse(time.DateTime, time.DateTime)
	return example
}

func (o *DateType) EmitToString(ref string, _ *GeneratorContext) string {
	ref = strings.TrimPrefix(ref, "*")
	return fmt.Sprintf("%s.Format(time.RFC3339)", ref)
}

func (o *DateType) IsPointerType() bool {
	return false
}
