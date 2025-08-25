package generator

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
)

var _ SchemaType = &IntType{}

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

func (o *IntType) BuildExample(*GeneratorContext, int, int) any {
	if ex := o.schema.Schema().Example; ex != nil {
		var decoded int64
		if err := ex.Decode(&decoded); err == nil {
			return decoded
		}
	}

	return 42
}

func (o *IntType) EmitToString(ref string, _ *GeneratorContext) string {
	return fmt.Sprintf("fmt.Sprintf(\"%%d\", %s)", ref)
}

func (o *IntType) IsPointerType() bool {
	return false
}
