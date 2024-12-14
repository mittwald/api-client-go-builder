package generator

import (
	"encoding/json"
	"github.com/moznion/gowrtr/generator"
)

var _ Type = &FloatType{}

type FloatType struct {
	BaseType
}

func (o *FloatType) IsLightweight() bool {
	return true
}

func (o *FloatType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	output, _ := json.Marshal(o.Schema.Schema())
	return []generator.Statement{
		generator.NewComment(string(output)),
		generator.NewRawStatementf("type %s = float64", o.Names.StructName),
	}
}

func (o *FloatType) EmitReference(*GeneratorContext) string {
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	return "float64"
}
