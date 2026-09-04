package generator

import (
	"github.com/moznion/gowrtr/generator"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

type Type interface {
	Name() SchemaName
	EmitDeclaration(ctx *GeneratorContext) []generator.Statement
	EmitReference(ctx *GeneratorContext) string
	IsPointerType() bool
}

type SchemaType interface {
	Type
	BuildExample(ctx *GeneratorContext, level, maxLevel int) any
	Schema() *base.SchemaProxy
	IsLightweight() bool
}

type TypeWithTestcases interface {
	EmitTestCases(ctx *GeneratorContext) []generator.Statement
}

type TypeWithValidation interface {
	EmitValidation(ref string, ctx *GeneratorContext) string
}

type TypeWithSubtypes interface {
	BuildSubtypes(opts GeneratorOpts, store *TypeStore) error
}

// TypeWithDeclarationName is implemented by types whose declaration is emitted
// under a different name than the one Name() reports. ReferenceType is the only
// such type: it adopts the name of its target so that references to it resolve
// to the target type, but a reference that is itself a component schema (an
// alias like `Foo: {$ref: Bar}`) still needs to be declared as `type Foo = Bar`
// under its own name and in its own file.
type TypeWithDeclarationName interface {
	DeclarationName() SchemaName
}

// declarationName returns the name a type's declaration should be emitted
// under, which is the type's own name unless it opts out via
// TypeWithDeclarationName.
func declarationName(typ Type) SchemaName {
	if d, ok := typ.(TypeWithDeclarationName); ok {
		return d.DeclarationName()
	}
	return typ.Name()
}
