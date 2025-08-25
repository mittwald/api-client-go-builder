package generator

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
)

var _ SchemaType = &FloatType{}

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

func (o *FloatType) BuildExample(*GeneratorContext, int, int) any {
	if ex := o.schema.Schema().Example; ex != nil {
		var decoded float64
		if err := ex.Decode(&decoded); err == nil {
			return decoded
		}
	}

	return 3.14
}

func (o *FloatType) EmitToString(ref string, _ *GeneratorContext) string {
	return fmt.Sprintf("fmt.Sprintf(\"%%f\", %s)", ref)
}

func (o *FloatType) IsPointerType() bool {
	return false
}
