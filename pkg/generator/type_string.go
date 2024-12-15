package generator

import (
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
	return []generator.Statement{
		generator.NewRawStatementf("type %s = string", o.Names.StructName),
	}
}

func (o *StringType) EmitReference(*GeneratorContext) string {
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	return "string"
}

func (o *StringType) BuildExample(*GeneratorContext) any {
	if example := o.schema.Schema().Example; example != nil {
		return example.Value
	}

	if examples := o.schema.Schema().Examples; len(examples) > 0 {
		return examples[0].Value
	}

	return "string"
}
