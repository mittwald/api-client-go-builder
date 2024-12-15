package generator

import (
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
	"path"
	"strings"
)

type SchemaName struct {
	PackageKey     string
	PackagePath    string
	StructName     string
	ForceNamedType bool
}

func (n *SchemaName) ForSubtype(subtype string) SchemaName {
	n2 := *n
	n2.StructName += util.UpperFirst(subtype)
	n2.PackagePath = strings.Replace(n2.PackagePath, ".go", "_"+strings.ToLower(subtype)+".go", 1)
	n2.ForceNamedType = false

	return n2
}

func (n *SchemaName) ForTestcase() SchemaName {
	n2 := *n
	n2.PackagePath = strings.Replace(n2.PackagePath, ".go", "_test.go", 1)
	n2.PackageKey += "_test"

	return n2
}

func (n *SchemaName) BuildRoot() *generator.Root {
	return generator.NewRoot(
		generator.NewComment(" THIS CODE WAS AUTO GENERATED"),
		generator.NewPackage(n.PackageKey),
		// we're running goimports after generating, so it does not matter if the uuid package is actually needed in a file
		generator.NewImport("github.com/google/uuid"),
		generator.NewNewline(),
	)
}

type SchemaNamingStrategy func(schemaName string) SchemaName

func MittwaldV1Strategy(schemaName string) SchemaName {
	// Example:
	// de.mittwald.v1.sshuser.SshUser
	parts := strings.Split(schemaName, ".")
	name := util.UpperFirst(parts[len(parts)-1])
	pkg := parts[len(parts)-2]
	version := parts[len(parts)-3]

	return SchemaName{
		StructName:     name,
		PackageKey:     pkg + version,
		PackagePath:    path.Join("schemas", pkg+version, strings.ToLower(name)+".go"),
		ForceNamedType: true,
	}

}
