package generator

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
)

var _ SchemaType = &ReferenceType{}

type ReferenceType struct {
	BaseType
	Target     string
	TargetType SchemaType
}

func (o *ReferenceType) IsLightweight() bool {
	return true
}

func (o *ReferenceType) BuildSubtypes(store *TypeStore) error {
	target, err := store.LookupReference(o.Target)
	if err != nil {
		return fmt.Errorf("could not resolve reference '%s': %w", o.Target, err)
	}

	o.BaseType.Names = target.Name()
	o.TargetType = target
	return nil
}

func (o *ReferenceType) lookupReferenceOnce(ctx *GeneratorContext) SchemaType {
	if o.TargetType != nil {
		return o.TargetType
	}

	o.TargetType, _ = ctx.KnownTypes.LookupReference(o.Target)
	return o.TargetType
}

func (o *ReferenceType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	return []generator.Statement{
		generator.NewRawStatementf("type %s = %s", o.Names.StructName, o.EmitReference(ctx)),
	}
}

func (o *ReferenceType) EmitReference(ctx *GeneratorContext) string {
	target := o.lookupReferenceOnce(ctx)
	if target != nil {
		return target.EmitReference(ctx)
	}

	return fmt.Sprintf("ERROR /* could not resolve %s */", o.Target)
	//return "any"
}

func (o *ReferenceType) EmitValidation(ref string, ctx *GeneratorContext) string {
	target := o.lookupReferenceOnce(ctx)
	if v, ok := target.(TypeWithValidation); ok {
		return v.EmitValidation(ref, ctx)
	}
	return "nil"
}

func (o *ReferenceType) BuildExample(ctx *GeneratorContext, level, maxLevel int) any {
	target := o.lookupReferenceOnce(ctx)
	return target.BuildExample(ctx, level+1, maxLevel)
}
