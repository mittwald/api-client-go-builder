package generator

import (
	"encoding/json"
	"fmt"
	"github.com/moznion/gowrtr/generator"
	"strings"
)

var _ Type = &ArrayType{}

type ArrayType struct {
	BaseType
	ItemType Type
}

func (o *ArrayType) BuildSubtypes(store *TypeStore) error {
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
	output, _ := json.Marshal(o.Schema.Schema())
	return []generator.Statement{
		generator.NewComment(string(output)),
		generator.NewCommentf("item type: %#v", o.ItemType),
		generator.NewRawStatementf("type %s = []%s", o.Names.StructName, o.ItemType.EmitReference(ctx)),
	}
}

func (o *ArrayType) EmitReference(ctx *GeneratorContext) string {
	innerRef := o.ItemType.EmitReference(ctx)
	if o.ItemType.IsLightweight() || strings.HasSuffix(innerRef, "Item") {
		return "[]" + o.ItemType.EmitReference(ctx)
	}

	if ctx.CurrentPackage == o.Names.PackageKey {
		return o.Names.StructName
	}

	return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
}
