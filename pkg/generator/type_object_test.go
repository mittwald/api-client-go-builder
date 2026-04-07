package generator

import (
	"testing"
)

func TestBuildSubtypesKeepsNestedInlineObjectProperties(t *testing.T) {
	spec := []byte(`{
  "openapi": "3.0.3",
  "info": {
    "title": "test",
    "version": "1.0.0"
  },
  "paths": {},
  "components": {
    "schemas": {
      "de.mittwald.v1.example.Parent": {
        "type": "object",
        "required": ["items"],
        "properties": {
          "items": {
            "type": "array",
            "items": {
              "type": "object",
              "required": ["child"],
              "properties": {
                "child": {
                  "type": "object",
                  "required": ["name"],
                  "properties": {
                    "name": {"type": "string"}
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`)

	doc, err := buildSpec(spec)
	if err != nil {
		t.Fatalf("buildSpec failed: %v", err)
	}

	store := NewTypeStore()
	sg := SchemaGenerator{SchemaNamingStrategy: MittwaldAPIVersionSchemaStrategy("v2")}

	for schemaName, schema := range doc.Model.Components.Schemas.FromOldest() {
		typ, err := sg.Build(schemaName, schema, store)
		if err != nil {
			t.Fatalf("building schema %q failed: %v", schemaName, err)
		}
		store.AddComponentSchema(schemaName, typ)
	}

	if err := store.BuildSubtypes(GeneratorOpts{}); err != nil {
		t.Fatalf("building subtypes failed: %v", err)
	}

	parentType, ok := store.ComponentSchemas["de.mittwald.v1.example.Parent"].(*ObjectType)
	if !ok {
		t.Fatalf("expected Parent to be ObjectType, got %T", store.ComponentSchemas["de.mittwald.v1.example.Parent"])
	}

	itemsType, _ := parentType.PropertyTypes.Get("items")
	itemsArray, ok := itemsType.(*ArrayType)
	if !ok {
		t.Fatalf("expected items to be ArrayType, got %T", itemsType)
	}

	itemObject, ok := itemsArray.ItemType.(*ObjectType)
	if !ok {
		t.Fatalf("expected items item type to be ObjectType, got %T", itemsArray.ItemType)
	}

	childType, _ := itemObject.PropertyTypes.Get("child")
	childObject, ok := childType.(*ObjectType)
	if !ok {
		t.Fatalf("expected child to be ObjectType, got %T", childType)
	}

	propertyCount := 0
	for range childObject.PropertyTypes.FromOldest() {
		propertyCount++
	}
	if propertyCount != 1 {
		t.Fatalf("expected child object to contain 1 property, got %d", propertyCount)
	}

	example := itemObject.BuildExample(&GeneratorContext{KnownTypes: store}, 0, 5)
	exampleMap, ok := example.(map[string]any)
	if !ok {
		t.Fatalf("expected example to be map[string]any, got %T", example)
	}

	childExample, ok := exampleMap["child"].(map[string]any)
	if !ok {
		t.Fatalf("expected child example to be map[string]any, got %T", exampleMap["child"])
	}

	if _, found := childExample["name"]; !found {
		t.Fatalf("expected child example to contain name, got %#v", childExample)
	}
}
