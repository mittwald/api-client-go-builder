package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
	"github.com/pb33f/libopenapi/orderedmap"
)

var _ Type = &ObjectType{}

type ObjectType struct {
	BaseType

	PropertyTypes      *orderedmap.Map[string, Type]
	RequiredProperties map[string]struct{}
}

func (o *ObjectType) IsLightweight() bool {
	return false
}

func (o *ObjectType) BuildSubtypes(store *TypeStore) error {
	s := o.schema.Schema()

	o.PropertyTypes = orderedmap.New[string, Type]()
	o.RequiredProperties = make(map[string]struct{})
	for _, req := range s.Required {
		o.RequiredProperties[req] = struct{}{}
	}

	for propName, propSchema := range s.Properties.FromOldest() {
		subTypeName := o.Names.ForSubtype(propName)

		propertyType, err := BuildTypeFromSchema(subTypeName, propSchema, store)
		if err != nil {
			return fmt.Errorf("error building subtype for property '%s': %w", propName, err)
		}

		if _, isRequired := o.RequiredProperties[propName]; !isRequired {
			propertyType = &OptionalType{InnerType: propertyType}
		}

		store.AddSubtype(subTypeName, propertyType)
		o.PropertyTypes.Set(propName, propertyType)
	}
	return nil
}

func (o *ObjectType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	s := o.schema.Schema()

	structDecl := generator.NewStruct(o.Names.StructName)
	for propName, propType := range o.PropertyTypes.FromOldest() {
		jsonTag := propName
		if _, isRequired := o.RequiredProperties[propName]; !isRequired {
			jsonTag += ",omitempty"
		}

		fieldName := util.UpperFirst(propName)
		structDecl = structDecl.AddField(
			fieldName,
			propType.EmitReference(ctx),
			fmt.Sprintf("json:\"%s\"", jsonTag),
		)
	}

	return []generator.Statement{
		generator.NewComment(s.Description),
		structDecl,
	}
}

func (o *ObjectType) EmitReference(ctx *GeneratorContext) string {
	if ctx.CurrentPackage == o.Names.PackageKey {
		return o.Names.StructName
	}
	return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
}
