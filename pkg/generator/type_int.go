package generator

import (
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
	return []generator.Statement{
		generator.NewRawStatementf("type %s = int64", o.Names.StructName),
	}
}

func (o *IntType) EmitReference(*GeneratorContext) string {
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	return "int64"
}
