package generator

import "github.com/moznion/gowrtr/generator"

type Type interface {
	Name() SchemaName
	IsLightweight() bool
	EmitDeclaration(ctx *GeneratorContext) []generator.Statement
	EmitReference(ctx *GeneratorContext) string
}

type TypeWithSubtypes interface {
	BuildSubtypes(store *TypeStore) error
}
