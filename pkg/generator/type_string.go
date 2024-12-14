package generator

import (
	"encoding/json"
	"github.com/moznion/gowrtr/generator"
)

var _ Type = &StringType{}

type StringType struct {
	BaseType
}

func (o *StringType) IsLightweight() bool {
	return true
}

func (o *StringType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	output, _ := json.Marshal(o.Schema.Schema())
	return []generator.Statement{
		generator.NewComment(string(output)),
		generator.NewRawStatementf("type %s = string", o.Names.StructName),
	}
}

func (o *StringType) EmitReference(*GeneratorContext) string {
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	return "string"
}
