package generator

import (
	"fmt"
	"github.com/charmbracelet/log"
)

type Generator struct {
	SpecLoader      SpecLoader
	SchemaGenerator SchemaGenerator
}

type GeneratorOpts struct {
	SpecSource      string
	Target          string
	BasePackageName string
}

func (g *Generator) Build(opts GeneratorOpts) error {
	log.Info("loading spec", "source", opts.SpecSource)

	doc, err := g.SpecLoader.LoadSpec(opts.SpecSource)
	if err != nil {
		return err
	}

	store := NewTypeStore()

	log.Info("processing #/components/schemas...")
	for schemaName, schema := range doc.Model.Components.Schemas.FromOldest() {
		typ, err := g.SchemaGenerator.Build(schemaName, schema, store)
		if err != nil {
			return fmt.Errorf("error generating schema '%s': %w", schemaName, err)
		}

		log.Debug("observed type", "name", schemaName)

		store.AddComponentSchema(schemaName, typ)
	}

	if err := g.generateClients(opts, doc, store); err != nil {
		return fmt.Errorf("error building clients: %w", err)
	}

	if err := store.BuildSubtypes(); err != nil {
		return fmt.Errorf("error while building subtypes: %w", err)
	}

	if err := store.EmitDeclarations(opts.Target); err != nil {
		return fmt.Errorf("error while emitting types: %w", err)
	}

	return nil
}
