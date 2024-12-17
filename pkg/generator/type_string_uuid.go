package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/moznion/gowrtr/generator"
)

var _ SchemaType = &StringUUIDType{}

type StringUUIDType struct {
	BaseType
}

func (o *StringUUIDType) IsLightweight() bool {
	return true
}

func (o *StringUUIDType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	stmts := make([]generator.Statement, 0)

	if d := o.schema.Schema().Description; d != "" {
		stmts = append(stmts, generatorx.NewMultilineComment(d))
	}

	stmts = append(stmts, generator.NewRawStatementf("type %s uuid.UUID", o.Names.StructName))
	return stmts
}

func (o *StringUUIDType) EmitReference(ctx *GeneratorContext) string {
	return "uuid.UUID"
}

func (o *StringUUIDType) BuildExample(*GeneratorContext, int, int) any {
	if example := o.schema.Schema().Example; example != nil {
		return example.Value
	}

	if examples := o.schema.Schema().Examples; len(examples) > 0 {
		return examples[0].Value
	}

	return "7a9d8971-09b0-4c39-8c64-546b6e1875ce"
}

func (o *StringUUIDType) EmitToString(ref string, _ *GeneratorContext) string {
	return fmt.Sprintf("%s.String()", ref)
}
