package generator

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/moznion/gowrtr/generator"
)

var _ Type = &ReferenceType{}

type ReferenceType struct {
	BaseType
	Target     string
	TargetType Type
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

func (o *ReferenceType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	return []generator.Statement{
		generator.NewRawStatementf("type %s = %s", o.Names.StructName, o.EmitReference(ctx)),
	}
}

func (o *ReferenceType) EmitReference(ctx *GeneratorContext) string {
	target, err := ctx.KnownTypes.LookupReference(o.Target)
	if err == nil {
		return target.EmitReference(ctx)
	}

	log.Warn("could not resolve reference", "ref", o.Target, "err", err)
	return fmt.Sprintf("ERROR /* could not resolve %s */", o.Target)
	//return "any"
}

func (o *ReferenceType) EmitValidation(ref string, ctx *GeneratorContext) string {
	if v, ok := o.TargetType.(TypeWithValidation); ok {
		return v.EmitValidation(ref, ctx)
	}
	return "nil"
}

func (o *ReferenceType) BuildExample() any {
	return o.TargetType.BuildExample()
}
