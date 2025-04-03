package generator

import (
	"fmt"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func BuildTypeFromSchema(names SchemaName, schema *base.SchemaProxy, knownTypes *TypeStore) (SchemaType, error) {
	baseType := BaseType{Names: names, schema: schema}
	format := schema.Schema().Format

	if schema.IsReference() {
		return &ReferenceType{BaseType: baseType, Target: schema.GetReference()}, nil
	}

	if len(schema.Schema().OneOf) > 0 {
		alternativeTypes := make([]SchemaType, len(schema.Schema().OneOf))
		for i, altSchema := range schema.Schema().OneOf {
			altType, err := BuildTypeFromSchema(names.ForSubtype(fmt.Sprintf("alternative%d", i+1)), altSchema, knownTypes)
			if err != nil {
				return nil, fmt.Errorf("error building alternative type %d for %s: %w", i, names.StructName, err)
			}

			alternativeTypes[i] = altType
		}

		return &OneOfType{BaseType: baseType, AlternativeTypes: alternativeTypes}, nil
	}

	// This is a hack used in some of our API routes to declare fields as explicitly nullable
	// (as opposed to "optional"), for example for PATCH routes in which "null" has a different
	// semantic than "not set".
	if len(schema.Schema().AllOf) == 1 && schema.Schema().Nullable != nil && *schema.Schema().Nullable {
		innerType, err := BuildTypeFromSchema(names, schema.Schema().AllOf[0], knownTypes)
		if err != nil {
			return nil, fmt.Errorf("error building inner type for %s: %w", names.StructName, err)
		}

		return &ExplicitlyNullableType{
			BaseType:  baseType,
			InnerType: innerType,
		}, nil
	}

	schemaType, err := GuessTypeFromSchema(schema)
	if err != nil {
		return nil, err
	}

	switch schemaType {
	case "object":
		additionalProperties := schema.Schema().AdditionalProperties
		if additionalProperties != nil && (additionalProperties.IsA() || (additionalProperties.IsB() && additionalProperties.B)) {
			var itemType SchemaType = &UnknownType{baseType}
			var err error

			if additionalProperties.IsA() {
				itemType, err = BuildTypeFromSchema(names.ForSubtype("item"), additionalProperties.A, knownTypes)
			}

			if err != nil {
				return nil, fmt.Errorf("error building array item type for %s: %w", names.StructName, err)
			}

			return &MapType{BaseType: baseType, ItemType: itemType}, nil
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
		switch format {
		case "uuid":
			return &StringUUIDType{BaseType: baseType}, nil
		case "date-time":
			return &DateType{BaseType: baseType}, nil
		default:
			return &StringType{BaseType: baseType}, nil
		}
	case "bool", "boolean":
		return &BoolType{BaseType: baseType}, nil
	case "integer":
		return &IntType{BaseType: baseType}, nil
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
	if schema.Schema().AdditionalProperties != nil {
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
