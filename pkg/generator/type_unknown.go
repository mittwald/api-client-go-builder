package generator

import (
	"encoding/json"
	"github.com/moznion/gowrtr/generator"
)

var _ Type = &UnknownType{}

type UnknownType struct {
	BaseType
}

func (o *UnknownType) IsLightweight() bool {
	return true
}

func (o *UnknownType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	output, _ := json.Marshal(o.schema.Schema())
	return []generator.Statement{
		generator.NewCommentf("TODO: This schema could not be automatically generated"),
		generator.NewComment(string(output)),
		generator.NewRawStatementf("type %s = any", o.Names.StructName),
	}
}

func (o *UnknownType) EmitReference(*GeneratorContext) string {
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	return "any"
}

func (o *UnknownType) BuildExample(*GeneratorContext, int, int) any {
	return nil
}
