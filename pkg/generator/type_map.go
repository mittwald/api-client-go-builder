package generator

import (
	"github.com/moznion/gowrtr/generator"
	"strings"
)

var _ SchemaType = &MapType{}

type MapType struct {
	BaseType
	ItemType SchemaType
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
	return []generator.Statement{
		generator.NewCommentf("item type: %#v", o.ItemType),
		generator.NewRawStatementf("type %s = map[string]%s", o.Names.StructName, o.ItemType.EmitReference(ctx)),
	}
}

func (o *MapType) EmitReference(ctx *GeneratorContext) string {
	return "map[string]" + o.ItemType.EmitReference(ctx)
}

func (o *MapType) BuildExample(ctx *GeneratorContext, level, maxLevel int) any {
	if level == maxLevel {
		return map[string]any{}
	}

	if ex := o.schema.Schema().Example; ex != nil {
		var decoded map[string]any
		if err := ex.Decode(&decoded); err == nil {
			return decoded
		}
	}

	return map[string]any{
		"string": o.ItemType.BuildExample(ctx, level+1, maxLevel),
	}
}
