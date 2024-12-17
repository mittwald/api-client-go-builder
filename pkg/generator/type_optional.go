package generator

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

var _ SchemaType = &OptionalType{}

type OptionalType struct {
	InnerType SchemaType
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

func (o *OptionalType) BuildSubtypes(store *TypeStore) error {
	if s, ok := o.InnerType.(TypeWithSubtypes); ok {
		return s.BuildSubtypes(store)
	}
	return nil
}

func (o *OptionalType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	return o.InnerType.EmitDeclaration(ctx)
}

func (o *OptionalType) EmitReference(ctx *GeneratorContext) string {
	innerType := o.InnerType
	if ref, isReference := o.InnerType.(*ReferenceType); isReference {
		innerType, _ = ctx.KnownTypes.LookupReference(ref.Target)
	}

	// slices are nil-able, anyway
	if _, isSlice := innerType.(*ArrayType); isSlice {
		return innerType.EmitReference(ctx)
	}

	if _, isMap := innerType.(*MapType); isMap {
		return innerType.EmitReference(ctx)
	}

	return "*" + innerType.EmitReference(ctx)
}

func (o *OptionalType) EmitValidation(ref string, ctx *GeneratorContext) string {
	if v, ok := o.InnerType.(TypeWithValidation); ok {
		return fmt.Sprintf("func () error {\nif %s == nil {\nreturn nil\n}\nreturn %s\n}()", ref, v.EmitValidation(ref, ctx))
	}
	return "nil"
}

func (o *OptionalType) BuildExample(ctx *GeneratorContext, level, maxLevel int) any {
	if level == maxLevel {
		return nil
	}

	return o.InnerType.BuildExample(ctx, level+1, maxLevel)
}

func (o *OptionalType) EmitToString(ref string, ctx *GeneratorContext) string {
	if ts, ok := o.InnerType.(TypeWithStringConversion); ok {
		return ts.EmitToString("*"+ref, ctx)
	}

	// if they want compile errors, give them compile errors!
	return "invalid-no-string-conversion"
}
