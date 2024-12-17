package generator

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"path"
	"strings"
)

var _ Type = &ClientSet{}

type ClientSet struct {
	name    SchemaName
	spec    *v3.Document
	clients *orderedmap.OrderedMap[string, *Client]
}

func (c *ClientSet) Name() SchemaName {
	return c.name
}

func (c *ClientSet) collectOperationsWithTag(tag string) []OperationWithMeta {
	operations := make([]OperationWithMeta, 0)

	for urlPath, items := range c.spec.Paths.PathItems.FromOldest() {
		for method, op := range items.GetOperations().FromOldest() {
			if util.SliceContains(op.Tags, tag) {
				operations = append(operations, OperationWithMeta{
					Path:      urlPath,
					Method:    method,
					Operation: op,
				})
			}
		}
	}

	return operations
}

func (c *ClientSet) BuildSubtypes(store *TypeStore) error {
	c.clients = orderedmap.New[string, *Client]()

	for _, tag := range c.spec.Tags {
		clientFunctionName := util.ConvertToTypename(tag.Name)
		clientTypeName := "Client"

		clientNameSet := c.name
		clientNameSet.StructName = clientTypeName
		clientNameSet.PackageKey = strings.ToLower(clientFunctionName)
		clientNameSet.PackagePath = path.Join(path.Dir(clientNameSet.PackagePath), strings.ToLower(clientFunctionName), "client_interface.go")

		client := Client{
			name:       clientNameSet,
			operations: c.collectOperationsWithTag(tag.Name),
		}

		log.Info("building client", "name", clientNameSet, "opcount", len(client.operations))

		store.AddClient(&client)
		c.clients.Set(clientFunctionName, &client)
	}

	return nil
}

func (c *ClientSet) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	iface := generator.NewInterface(c.name.StructName)

	for clientName, client := range c.clients.FromOldest() {
		iface = iface.AddSignatures(generator.NewFuncSignature(clientName).AddReturnTypes(client.EmitReference(ctx)))
	}

	return []generator.Statement{iface}
}

func (c *ClientSet) EmitReference(ctx *GeneratorContext) string {
	if c.name.PackageKey == ctx.CurrentPackage {
		return c.name.StructName
	}

	return fmt.Sprintf("%s.%s", c.name.PackageKey, c.name.StructName)
}
