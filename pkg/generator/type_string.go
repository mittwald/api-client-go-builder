package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/moznion/gowrtr/generator"
)

var _ SchemaType = &StringType{}

type StringType struct {
	BaseType
}

func (o *StringType) IsLightweight() bool {
	return true
}

func (o *StringType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	stmts := make([]generator.Statement, 0)

	if d := o.schema.Schema().Description; d != "" {
		stmts = append(stmts, generatorx.NewMultilineComment(d))
	}

	stmts = append(stmts, generator.NewRawStatementf("type %s string", o.Names.StructName))
	return stmts
}

func (o *StringType) EmitReference(ctx *GeneratorContext) string {
	if o.Names.ForceNamedType {
		if ctx.CurrentPackage == o.Names.PackageKey {
			return o.Names.StructName
		}
		return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	}
	return "string"
}

func (o *StringType) BuildExample(*GeneratorContext, int, int) any {
	if example := o.schema.Schema().Example; example != nil {
		return example.Value
	}

	if examples := o.schema.Schema().Examples; len(examples) > 0 {
		return examples[0].Value
	}

	return "string"
}

func (o *StringType) EmitToString(ref string, _ *GeneratorContext) string {
	return ref
}
