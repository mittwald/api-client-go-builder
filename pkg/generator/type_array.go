package generator

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
	"strings"
)

var _ SchemaType = &ArrayType{}

type ArrayType struct {
	BaseType
	ItemType SchemaType
}

func (o *ArrayType) BuildSubtypes(opts GeneratorOpts, store *TypeStore) error {
	if s, ok := o.ItemType.(TypeWithSubtypes); ok {
		if err := s.BuildSubtypes(opts, store); err != nil {
			return err
		}
	}

	if o.ItemType.IsLightweight() {
		return nil
	}

	subTypeName := o.Names.ForSubtype("item")

	store.AddSubtype(subTypeName, o.ItemType)
	return nil
}

func (o *ArrayType) IsLightweight() bool {
	return o.ItemType.IsLightweight() || strings.HasSuffix(o.ItemType.Name().StructName, "Item")
}

func (o *ArrayType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	return []generator.Statement{
		generator.NewRawStatementf("type %s []%s", o.Names.StructName, o.ItemType.EmitReference(ctx)),
	}
}

func (o *ArrayType) EmitReference(ctx *GeneratorContext) string {
	if o.Names.ForceNamedType {
		if ctx.CurrentPackage == o.Names.PackageKey {
			return o.Names.StructName
		}

		return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	}

	innerRef := o.ItemType.EmitReference(ctx)
	if o.ItemType.IsLightweight() || strings.HasSuffix(innerRef, "Item") {
		return "[]" + o.ItemType.EmitReference(ctx)
	}

	if ctx.CurrentPackage == o.Names.PackageKey {
		return o.Names.StructName
	}

	return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
}

func (o *ArrayType) EmitValidation(ref string, ctx *GeneratorContext) string {
	if v, ok := o.ItemType.(TypeWithValidation); ok {
		validation := v.EmitValidation(ref+"[i]", ctx)
		if validation == "nil" {
			return validation
		}

		return fmt.Sprintf("func () error {\n"+
			"for i := range %s {\n"+
			"if err := %s; err != nil {\n"+
			"return fmt.Errorf(\"item %%d is invalid %%w\", i, err)\n"+
			"}\n"+
			"}\n"+
			"return nil\n"+
			"}()", ref, v.EmitValidation(ref+"[i]", ctx))
	}
	return "nil"
}

func (o *ArrayType) BuildExample(ctx *GeneratorContext, level, maxLevel int) any {
	if level == maxLevel {
		return []any{}
	}
	return []any{o.ItemType.BuildExample(ctx, level+1, maxLevel)}
}
