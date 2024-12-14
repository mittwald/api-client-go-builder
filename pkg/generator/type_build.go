package generator

import (
	"fmt"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func BuildTypeFromSchema(names SchemaName, schema *base.SchemaProxy, knownTypes *TypeStore) (Type, error) {
	baseType := BaseType{Names: names, Schema: schema}
	format := schema.Schema().Format

	if schema.IsReference() {
		return &ReferenceType{BaseType: baseType, Target: schema.GetReference()}, nil
		//ref := schema.GetReference()
		//return knownTypes.LookupReference(ref)
	}

	schemaType, err := GuessTypeFromSchema(schema)
	if err != nil {
		return nil, err
	}

	switch schemaType {
	case "object":
		if schema.Schema().AdditionalProperties != nil {
			// TODO: Should be mapped to a map type, instead
			return &UnknownType{BaseType: baseType}, nil
		}
		return &ObjectType{BaseType: baseType}, nil
	case "array":
		items := schema.Schema().Items
		if items != nil && items.IsA() {
			itemType, err := BuildTypeFromSchema(names.ForSubtype("item"), items.A, knownTypes)
			if err != nil {
				return nil, fmt.Errorf("error building array item type for %s: %w", names.StructName, err)
			}
			return &ArrayType{BaseType: baseType, ItemType: itemType}, nil
		}

		return &ArrayType{BaseType: baseType, ItemType: &UnknownType{BaseType: baseType}}, nil
	case "string":
		if schema.Schema().Enum != nil {
			return NewStringEnumTypeFromYamlNodes(baseType, schema.Schema().Enum), nil
		}
		return &StringType{BaseType: baseType}, nil
	case "bool", "boolean":
		return &BoolType{BaseType: baseType}, nil
	case "number":
		if format == "int" || format == "integer" {
			return &IntType{BaseType: baseType}, nil
		}
		return &FloatType{BaseType: baseType}, nil
	default:
		return &UnknownType{BaseType: baseType}, nil
	}
}

func GuessTypeFromSchema(schema *base.SchemaProxy) (string, error) {
	schemaTypes := schema.Schema().Type

	if len(schemaTypes) == 1 {
		return schemaTypes[0], nil
	}

	// When there are properties, assume the schema to be an object
	if schema.Schema().Properties != nil {
		return "object", nil
	}

	if schema.Schema().Items != nil {
		return "array", nil
	}

	if len(schemaTypes) > 1 {
		// TODO: Support multiple types (necessary if the spec ever switches to OpenAPI 3.1)
		return "", fmt.Errorf("schemas with multiple types are not supported")
	}

	return "unknown", nil
}
