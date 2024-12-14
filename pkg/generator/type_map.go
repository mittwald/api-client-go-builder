package generator

import (
	"encoding/json"
	"github.com/moznion/gowrtr/generator"
	"strings"
)

var _ Type = &MapType{}

type MapType struct {
	BaseType
	ItemType Type
}

func (o *MapType) BuildSubtypes(store *TypeStore) error {
	if o.ItemType.IsLightweight() {
		return nil
	}

	subTypeName := o.Names.ForSubtype("item")

	store.AddSubtype(subTypeName, o.ItemType)
	return nil
}

func (o *MapType) IsLightweight() bool {
	return o.ItemType.IsLightweight() || strings.HasSuffix(o.ItemType.Name().StructName, "Item")
}

func (o *MapType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	output, _ := json.Marshal(o.Schema.Schema())
	return []generator.Statement{
		generator.NewComment(string(output)),
		generator.NewCommentf("item type: %#v", o.ItemType),
		generator.NewRawStatementf("type %s = map[string]%s", o.Names.StructName, o.ItemType.EmitReference(ctx)),
	}
}

func (o *MapType) EmitReference(ctx *GeneratorContext) string {
	return "map[string]" + o.ItemType.EmitReference(ctx)

	//innerRef := o.ItemType.EmitReference(ctx)
	//if o.ItemType.IsLightweight() || strings.HasSuffix(innerRef, "Item") {
	//	return "map[string]" + o.ItemType.EmitReference(ctx)
	//}
	//
	//if ctx.CurrentPackage == o.Names.PackageKey {
	//	return o.Names.StructName
	//}
	//
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
}
