package generator

import (
	"github.com/moznion/gowrtr/generator"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

type Type interface {
	Name() SchemaName
	EmitDeclaration(ctx *GeneratorContext) []generator.Statement
	EmitReference(ctx *GeneratorContext) string
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
	BuildSubtypes(store *TypeStore) error
}
