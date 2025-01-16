package generator

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/moznion/gowrtr/generator"
	"path"
	"strings"
)

type TypeStore struct {
	ComponentSchemas map[string]Type
	SubTypes         map[string]Type
	Clients          map[string]Type
}

func NewTypeStore() *TypeStore {
	return &TypeStore{
		ComponentSchemas: make(map[string]Type),
		SubTypes:         make(map[string]Type),
		Clients:          make(map[string]Type),
	}
}

func (s *TypeStore) LookupReference(ref string) (SchemaType, error) {
	if strings.HasPrefix(ref, "#/components/schemas") {
		name := strings.Replace(ref, "#/components/schemas/", "", 1)
		typ, ok := s.ComponentSchemas[name]
		if !ok {
			return nil, fmt.Errorf("reference '%s' could not be resolved; no type with this name", ref)
		}

		if schemaType, ok := typ.(SchemaType); ok {
			return schemaType, nil
		}
	}

	return nil, fmt.Errorf("reference '%s' could not be resolved; no type with this name", ref)
}

func (s *TypeStore) AddComponentSchema(name string, typ Type) {
	s.ComponentSchemas[name] = typ
}

func (s *TypeStore) AddClient(typ Type) {
	s.Clients[typ.Name().PackagePath] = typ
}

func (s *TypeStore) AddSubtype(name SchemaName, typ Type) {
	s.SubTypes[name.PackageKey+"."+name.StructName] = typ
}

func (s *TypeStore) Len() int {
	return len(s.ComponentSchemas) + len(s.SubTypes) + len(s.Clients)
}

func (s *TypeStore) BuildSubtypes() error {
	log.Info("building subtypes", "count", s.Len())

	visited := make(map[string]struct{})

	buildSubtype := func(name string, typ Type) error {
		if _, alreadySeen := visited[name]; alreadySeen {
			return nil
		}

		visited[name] = struct{}{}
		log.Debug("building subtypes for", "name", name)
		if st, ok := typ.(TypeWithSubtypes); ok {
			if err := st.BuildSubtypes(s); err != nil {
				return fmt.Errorf("error building subtypes for %s: %w", name, err)
			}
		}
		return nil
	}

	processAll := func() error {
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
		for name, typ := range s.Clients {
			if err := buildSubtype(name, typ); err != nil {
				return err
			}
		}
		return nil
	}

	processed := 0

	for processed < s.Len() {
		if err := processAll(); err != nil {
			return err
		}
		processed = len(visited)
	}

	if err := processAll(); err != nil {
		return err
	}

	return nil
}

func (s *TypeStore) EmitDeclarations(targetPath string) error {
	log.Info("emitting declarations", "count", s.Len())

	ctx := GeneratorContext{KnownTypes: s, WithDebuggingComments: true}

	packagesWithTestcases := make(map[string]string)

	buildType := func(what string, typ Type) error {
		ctxForType := ctx
		ctxForType.CurrentPackage = typ.Name().PackageKey

		names := typ.Name()
		stmts := typ.EmitDeclaration(&ctxForType)

		log.Infof("emitting declaration for %s", names.StructName)

		// without Goimports(); does not make sense until the very end
		root := names.BuildRoot()

		if ctxForType.WithDebuggingComments {
			if schemaType, ok := typ.(SchemaType); ok {
				schemaJson, _ := schemaType.Schema().Render()
				root = root.AddStatements(
					generator.NewComment("This data type was generated from the following JSON schema:"),
					generatorx.NewMultilineComment(strings.TrimRight(string(schemaJson), "\n")),
					generator.NewNewline(),
				)
			}
		}

		root = root.AddStatements(stmts...)
		//root = root.Gofmt("-s")

		if err := EmitToFile(targetPath, names.PackagePath, root); err != nil {
			return fmt.Errorf("error emitting %s %T to %s: %w", what, typ, names.PackagePath, err)
		}

		if tc, ok := typ.(TypeWithTestcases); ok {
			testNames := names.ForTestcase()
			testRoot := testNames.BuildRoot()

			packagesWithTestcases[testNames.PackageKey] = path.Join(path.Dir(names.PackagePath), "suite_test.go")

			testStmts := tc.EmitTestCases(&ctx)
			testRoot = testRoot.AddStatements(generator.NewRawStatement("import . \"github.com/onsi/ginkgo/v2\""))
			testRoot = testRoot.AddStatements(generator.NewRawStatement("import . \"github.com/onsi/gomega\""))
			testRoot = testRoot.AddStatements(testStmts...)
			//testRoot = testRoot.Gofmt("-s")

			if err := EmitToFile(targetPath, testNames.PackagePath, testRoot); err != nil {
				return err
			}
		}

		return nil
	}

	for _, typ := range s.ComponentSchemas {
		if err := buildType("schema", typ); err != nil {
			return err
		}
	}

	for _, typ := range s.Clients {
		if st, ok := typ.(SchemaType); ok {
			if st.IsLightweight() {
				continue
			}
		}

		if err := buildType("client", typ); err != nil {
			return err
		}
	}

	for _, typ := range s.SubTypes {
		if st, ok := typ.(SchemaType); ok {
			if st.IsLightweight() {
				log.Infof("skipping lightweight type %s (type %T)", st.Name().StructName, typ)
				continue
			}
		}

		if err := buildType("subtype", typ); err != nil {
			return err
		}
	}

	log.Info("generating test suites")
	for pkg, pkgPath := range packagesWithTestcases {
		suiteRoot := generator.NewRoot(
			generator.NewPackage(pkg),
			generator.NewNewline(),
			generator.NewRawStatement("import . \"github.com/onsi/ginkgo/v2\""),
			generator.NewRawStatement("import . \"github.com/onsi/gomega\""),
			generator.NewFunc(nil, generator.NewFuncSignature("TestTypes").AddParameters(generator.NewFuncParameter("t", "*testing.T")),
				generator.NewRawStatement("RegisterFailHandler(Fail)"),
				generator.NewRawStatementf("RunSpecs(t, \"%s types\")", pkg),
			),
		)

		if err := EmitToFile(targetPath, pkgPath, suiteRoot); err != nil {
			return err
		}
	}

	/*
		log.Info("running goimports")
		//cmd := exec.Command("goimports", "-w", targetPath)
		cmd := exec.Command("goimports", "-w", ".")
		cmd.Dir = targetPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error running goimports: %w", err)
		}

	*/

	return nil
}
