package generator

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

type SchemaGenerator struct {
	Opts                 GeneratorOpts
	SchemaNamingStrategy SchemaNamingStrategy
}

func (g *SchemaGenerator) Build(name string, schema *base.SchemaProxy, types *TypeStore) (Type, error) {
	names := g.SchemaNamingStrategy(name)

	typ, err := BuildTypeFromSchema(names, schema, types)
	if err != nil {
		return nil, err
	}

	return typ, nil
}
