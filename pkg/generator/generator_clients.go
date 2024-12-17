package generator

import (
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func (g *Generator) generateClients(opts GeneratorOpts, spec *libopenapi.DocumentModel[v3.Document], store *TypeStore) error {
	baseName := SchemaName{
		PackagePath: "clients/clientset.go",
		PackageKey:  opts.BasePackageName,
		StructName:  "Client",
	}

	clientSet := ClientSet{name: baseName, spec: &spec.Model}
	store.AddClient(&clientSet)

	return nil
}
