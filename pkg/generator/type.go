package generator

import (
	"github.com/moznion/gowrtr/generator"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

type Type interface {
	Name() SchemaName
	Schema() *base.SchemaProxy
	IsLightweight() bool
	EmitDeclaration(ctx *GeneratorContext) []generator.Statement
	EmitReference(ctx *GeneratorContext) string
	BuildExample() any
}

type TypeWithTestcases interface {
	EmitTestCases(ctx *GeneratorContext) []generator.Statement
}

type TypeWithValidation interface {
	EmitValidation(ref string, ctx *GeneratorContext) string
}

type TypeWithSubtypes interface {
	BuildSubtypes(store *TypeStore) error
}
