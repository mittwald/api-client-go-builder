package generator

import (
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
	return []generator.Statement{
		generator.NewRawStatementf("type %s = bool", o.Names.StructName),
	}
}

func (o *BoolType) EmitReference(*GeneratorContext) string {
	return "bool"
}

func (o *BoolType) BuildExample(*GeneratorContext) any {
	if ex := o.schema.Schema().Example; ex != nil {
		var decoded bool
		if err := ex.Decode(&decoded); err == nil {
			return decoded
		}
	}

	return true
}
