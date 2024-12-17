package generator

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"path"
	"strings"
)

type OperationWithMeta struct {
	Path      string
	Method    string
	Operation *v3.Operation

	requestType  Type
	responseType Type
}

type Client struct {
	name       SchemaName
	operations []OperationWithMeta
}

func (c *Client) Name() SchemaName {
	return c.name
}

func (c *Client) BuildSubtypes(store *TypeStore) error {
	for i, op := range c.operations {
		if op.Operation.OperationId == "" {
			log.Warn("empty operation id", "path", op.Path, "method", op.Method)
			continue
		}

		operationName := util.ConvertToTypename(op.Operation.OperationId)

		requestName := c.name
		requestName.StructName = operationName + "Request"
		requestName.PackagePath = path.Join(path.Dir(requestName.PackagePath), strings.ToLower(operationName)+"_request.go")

		responseName := c.name
		responseName.StructName = operationName + "Request"
		responseName.PackagePath = path.Join(path.Dir(requestName.PackagePath), strings.ToLower(operationName)+"_response.go")

		c.operations[i].requestType = &ClientOperationRequest{name: requestName, operation: &c.operations[i]}

		store.AddClient(c.operations[i].requestType)
	}

	return nil
}

func (c *Client) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	iface := generator.NewInterface(c.name.StructName)

	for _, op := range c.operations {
		if op.Operation.OperationId == "" {
			log.Warn("empty operation id", "path", op.Path, "method", op.Method)
			continue
		}

		operationName := util.ConvertToTypename(op.Operation.OperationId)
		request := operationName + "Request"
		response := operationName + "Response"

		iface = iface.AddSignatures(generator.NewFuncSignature(operationName).
			AddParameters(
				generator.NewFuncParameter("ctx", "context.Context"),
				generator.NewFuncParameter("req", request),
			).
			AddReturnTypes(response, "error"))
	}

	return []generator.Statement{iface}
}

func (c *Client) EmitReference(ctx *GeneratorContext) string {
	if c.name.PackageKey == ctx.CurrentPackage {
		return c.name.StructName
	}

	return fmt.Sprintf("%s.%s", c.name.PackageKey, c.name.StructName)
}
