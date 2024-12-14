package generator

import (
	"encoding/json"
	"github.com/moznion/gowrtr/generator"
)

var _ Type = &IntType{}

type IntType struct {
	BaseType
}

func (o *IntType) IsLightweight() bool {
	return true
}

func (o *IntType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	output, _ := json.Marshal(o.Schema.Schema())
	return []generator.Statement{
		generator.NewComment(string(output)),
		generator.NewRawStatementf("type %s = int64", o.Names.StructName),
	}
}

func (o *IntType) EmitReference(*GeneratorContext) string {
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	return "int64"
}
