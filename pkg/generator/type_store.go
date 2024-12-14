package generator

import (
	"fmt"
	"github.com/charmbracelet/log"
	"strings"
)

type TypeStore struct {
	ComponentSchemas map[string]Type
	SubTypes         map[string]Type
}

func NewTypeStore() *TypeStore {
	return &TypeStore{
		ComponentSchemas: make(map[string]Type),
		SubTypes:         make(map[string]Type),
	}
}

func (s *TypeStore) LookupReference(ref string) (Type, error) {
	if strings.HasPrefix(ref, "#/components/schemas") {
		name := strings.Replace(ref, "#/components/schemas/", "", 1)
		typ, ok := s.ComponentSchemas[name]
		if ok {
			return typ, nil
		}
	}

	return nil, fmt.Errorf("reference '%s' could not be resolved; no type with this name", ref)
}

func (s *TypeStore) AddComponentSchema(name string, typ Type) {
	s.ComponentSchemas[name] = typ
}

func (s *TypeStore) AddSubtype(name SchemaName, typ Type) {
	s.SubTypes[name.PackageKey+"."+name.StructName] = typ
}

func (s *TypeStore) Len() int {
	return len(s.ComponentSchemas) + len(s.SubTypes)
}

func (s *TypeStore) BuildSubtypes() error {
	log.Info("building subtypes", "count", s.Len())

	visited := make(map[string]struct{})

	buildSubtype := func(name string, typ Type) error {
		visited[name] = struct{}{}
		if st, ok := typ.(TypeWithSubtypes); ok {
			if err := st.BuildSubtypes(s); err != nil {
				return fmt.Errorf("error building subtypes for %s: %w", name, err)
			}
		}
		return nil
	}

	for len(visited) < s.Len() {
		for name, typ := range s.ComponentSchemas {
			if err := buildSubtype(name, typ); err != nil {
				return err
			}
		}
		for name, typ := range s.SubTypes {
			if err := buildSubtype(name, typ); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TypeStore) EmitDeclarations(targetPath string) error {
	log.Info("emitting declarations", "count", s.Len())

	ctx := GeneratorContext{KnownTypes: s}

	buildType := func(typ Type) error {
		ctxForType := ctx
		ctxForType.CurrentPackage = typ.Name().PackageKey

		names := typ.Name()
		stmts := typ.EmitDeclaration(&ctxForType)

		log.Infof("emitting declaration for %s", names.StructName)

		// without Goimports(); does not make sense until the very end
		root := names.BuildRoot().
			AddStatements(stmts...).
			Gofmt("-s")

		return EmitToFile(targetPath, names, root)
	}

	for _, typ := range s.ComponentSchemas {
		if err := buildType(typ); err != nil {
			return err
		}
	}

	for _, typ := range s.SubTypes {
		if typ.IsLightweight() {
			continue
		}

		if err := buildType(typ); err != nil {
			return err
		}
	}

	return nil
}
