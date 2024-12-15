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

func (o *IntType) BuildExample() any {
	if ex := o.schema.Schema().Example; ex != nil {
		var decoded int64
		if err := ex.Decode(&decoded); err == nil {
			return decoded
		}
	}

	return 42
}
