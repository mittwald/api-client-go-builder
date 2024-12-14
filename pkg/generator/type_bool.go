package generator

import (
	"encoding/json"
	"github.com/moznion/gowrtr/generator"
)

var _ Type = &BoolType{}

type BoolType struct {
	BaseType
}

func (o *BoolType) IsLightweight() bool {
	return true
}

func (o *BoolType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	output, _ := json.Marshal(o.Schema.Schema())
	return []generator.Statement{
		generator.NewComment(string(output)),
		generator.NewRawStatementf("type %s = bool", o.Names.StructName),
	}
}

func (o *BoolType) EmitReference(*GeneratorContext) string {
	return "bool"
}
