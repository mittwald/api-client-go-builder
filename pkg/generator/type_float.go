package generator

import (
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
	return []generator.Statement{
		generator.NewRawStatementf("type %s = float64", o.Names.StructName),
	}
}

func (o *FloatType) EmitReference(*GeneratorContext) string {
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	return "float64"
}
