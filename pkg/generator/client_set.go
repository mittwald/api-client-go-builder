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

func (c *ClientSet) BuildSubtypes(opts GeneratorOpts, store *TypeStore) error {
	c.clients = orderedmap.New[string, *Client]()

	for _, tag := range c.spec.Tags {
		clientFunctionName := util.ConvertToTypename(tag.Name)
		clientTypeName := "Client"
		clientPackageKey := strings.ToLower(clientFunctionName) + "client" + opts.APIVersion

		clientNameSet := c.name
		clientNameSet.StructName = clientTypeName
		clientNameSet.PackageKey = clientPackageKey
		clientNameSet.PackagePath = path.Join(path.Dir(clientNameSet.PackagePath), clientPackageKey, "client.go")

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

func (c *ClientSet) ImplName() string {
	return util.LowerFirst(fmt.Sprintf("%sImpl", c.name.StructName))
}

func (c *ClientSet) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	iface := generator.NewInterface(c.name.StructName)
	str := generator.NewStruct(c.ImplName()).AddField("client", "httpclient.RequestRunner")
	strReceiver := generator.NewFuncReceiver("c", "*"+c.ImplName())

	funcs := []generator.Statement{
		generator.NewFunc(
			nil,
			generator.NewFuncSignature("NewClient").
				Parameters(generator.NewFuncParameter("client", "httpclient.RequestRunner")).
				ReturnTypes(c.EmitReference(ctx)),
			generator.NewReturnStatement(fmt.Sprintf("&%s{client: client}", c.ImplName())),
		),
		generator.NewNewline(),
	}

	for clientName, client := range c.clients.FromOldest() {
		signature := generator.NewFuncSignature(clientName).AddReturnTypes(client.EmitReference(ctx))

		iface = iface.AddSignatures(signature)
		clientFunc := generator.NewFunc(strReceiver, signature, generator.NewReturnStatement(client.Name().PackageKey+".NewClient(c.client)"))

		funcs = append(funcs, clientFunc, generator.NewNewline())
	}

	stmts := []generator.Statement{iface, str}
	stmts = append(stmts, funcs...)

	return stmts
}

func (c *ClientSet) EmitReference(ctx *GeneratorContext) string {
	if c.name.PackageKey == ctx.CurrentPackage {
		return c.name.StructName
	}

	return fmt.Sprintf("%s.%s", c.name.PackageKey, c.name.StructName)
}
