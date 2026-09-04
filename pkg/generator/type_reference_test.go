package generator

import (
	"os"
	"path"
	"strings"
	"testing"
)

// A component schema may be a bare "$ref" to another component schema. Such an
// alias adopts its target's name for referencing purposes, but it must still be
// declared under its own name -- otherwise it overwrites (or, since EmitToFile
// refuses to overwrite, collides with) the file of its target.
func TestEmitDeclarationsForAliasedComponentSchemas(t *testing.T) {
	spec := []byte(`{
  "openapi": "3.0.3",
  "info": {
    "title": "test",
    "version": "1.0.0"
  },
  "paths": {},
  "components": {
    "schemas": {
      "de.mittwald.v1.example.Plan": {
        "type": "object",
        "required": ["id"],
        "properties": {
          "id": {"type": "string"}
        }
      },
      "de.mittwald.v1.example.PlanOptions": {
        "$ref": "#/components/schemas/de.mittwald.v1.example.Plan"
      },
      "de.mittwald.v1.example.CustomerPlanOptions": {
        "$ref": "#/components/schemas/de.mittwald.v1.example.Plan"
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

	targetPath := t.TempDir()
	if err := store.EmitDeclarations(targetPath, nil); err != nil {
		t.Fatalf("emitting declarations failed: %v", err)
	}

	for _, tc := range []struct {
		file    string
		content string
	}{
		{file: "schemas/examplev2/plan.go", content: "type Plan struct"},
		{file: "schemas/examplev2/planoptions.go", content: "type PlanOptions = Plan"},
		{file: "schemas/examplev2/customerplanoptions.go", content: "type CustomerPlanOptions = Plan"},
	} {
		emitted, err := os.ReadFile(path.Join(targetPath, tc.file))
		if err != nil {
			t.Fatalf("expected %s to be emitted: %v", tc.file, err)
		}

		if !strings.Contains(string(emitted), tc.content) {
			t.Fatalf("expected %s to contain %q, got:\n%s", tc.file, tc.content, emitted)
		}
	}
}
