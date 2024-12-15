package generator

import (
	"github.com/moznion/gowrtr/generator"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

var _ Type = &OptionalType{}

type OptionalType struct {
	InnerType Type
}

func (o *OptionalType) Name() SchemaName {
	return o.InnerType.Name()
}

func (o *OptionalType) Schema() *base.SchemaProxy {
	return o.InnerType.Schema()
}

func (o *OptionalType) IsLightweight() bool {
	return o.InnerType.IsLightweight()
}

func (o *OptionalType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	return o.InnerType.EmitDeclaration(ctx)
}

func (o *OptionalType) EmitReference(ctx *GeneratorContext) string {
	// slices are nil-able, anyway
	if _, isSlice := o.InnerType.(*ArrayType); isSlice {
		return o.InnerType.EmitReference(ctx)
	}

	if _, isMap := o.InnerType.(*MapType); isMap {
		return o.InnerType.EmitReference(ctx)
	}

	return "*" + o.InnerType.EmitReference(ctx)
}
