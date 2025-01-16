package generator

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"path"
	"strconv"
	"strings"
)

type OperationWithMeta struct {
	Path      string
	Method    string
	Operation *v3.Operation

	requestType    Type
	responseType   Type
	responseFormat string
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
		responseName.StructName = operationName + "Response"
		responseName.PackagePath = path.Join(path.Dir(requestName.PackagePath), strings.ToLower(operationName)+"_response.go")

		c.operations[i].requestType = &ClientOperationRequest{name: requestName, operation: &c.operations[i]}

		for code, response := range c.operations[i].Operation.Responses.Codes.FromOldest() {
			codeAsInt, err := strconv.ParseInt(code, 10, strconv.IntSize)
			if err != nil {
				return fmt.Errorf("response code %s of operation %s could not be parsed as int: %w", code, op.Operation.OperationId, err)
			}

			if codeAsInt >= 200 && codeAsInt < 400 && response.Content != nil {
				if schema, ok := response.Content.Get("application/json"); ok {
					responseType, err := BuildTypeFromSchema(responseName, schema.Schema, store)
					if err != nil {
						return fmt.Errorf("error building response type for operation %s: %w", op.Operation.OperationId, err)
					}

					c.operations[i].responseType = responseType
					c.operations[i].responseFormat = "json"

					store.AddClient(responseType)
				}
			}
		}

		store.AddClient(c.operations[i].requestType)
	}

	return nil
}

func (c *Client) ImplName() string {
	return util.LowerFirst(fmt.Sprintf("%sImpl", c.name.StructName))
}

func (c *Client) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	clientInterface := generator.NewInterface(c.name.StructName)
	clientStructName := c.ImplName()
	clientStruct := generator.NewStruct(clientStructName).
		AddField("client", "httpclient.RequestRunner")

	funcStmts := []generator.Statement{
		generator.NewFunc(
			nil,
			generator.NewFuncSignature("NewClient").
				Parameters(generator.NewFuncParameter("client", "httpclient.RequestRunner")).
				ReturnTypes(c.EmitReference(ctx)),
			generator.NewReturnStatement(fmt.Sprintf("&%s{client: client}", clientStructName)),
		),
	}

	clientStructReceiver := generator.NewFuncReceiver("c", "*"+clientStructName)

	for _, op := range c.operations {
		if op.Operation.OperationId == "" {
			log.Warn("empty operation id", "path", op.Path, "method", op.Method)
			continue
		}

		operationName := util.ConvertToTypename(op.Operation.OperationId)

		funcSignature := generator.NewFuncSignature(operationName).
			AddParameters(
				generator.NewFuncParameter("ctx", "context.Context"),
				generator.NewFuncParameter("req", op.requestType.EmitReference(ctx)),
			)

		errorReturn := generator.NewReturnStatement("nil", "err")
		errorReturnWithResponse := generator.NewReturnStatement("httpRes", "err")
		if op.responseType != nil {
			errorReturn = generator.NewReturnStatement("nil", "nil", "err")
			errorReturnWithResponse = generator.NewReturnStatement("nil", "httpRes", "err")
		}

		operationFuncStmts := []generator.Statement{
			generator.NewRawStatement("body, err := req.body()"),
			generator.NewIf("err != nil", errorReturn),
			generator.NewNewline(),
			generator.NewRawStatement("httpReq, err := http.NewRequestWithContext(ctx, req.method(), req.url(), body)"),
			generator.NewIf("err != nil", errorReturn),
			generator.NewNewline(),
			generator.NewRawStatement("httpRes, err := c.client.Do(httpReq)"),
			generator.NewIf("err != nil", errorReturnWithResponse),
			generator.NewNewline(),
		}

		if op.responseType != nil {
			funcSignature = funcSignature.AddReturnTypes("*"+op.responseType.EmitReference(ctx), "*http.Response", "error")
			if op.responseFormat == "json" {
				operationFuncStmts = append(operationFuncStmts,
					generator.NewRawStatement("decoder := json.NewDecoder(httpRes.Body)"),
					generator.NewRawStatementf("var response %s", op.responseType.EmitReference(ctx)),
					generator.NewIf("err := decoder.Decode(&response); err != nil", errorReturnWithResponse),
					generator.NewReturnStatement("&response", "httpRes", "nil"),
				)
			} else {
				operationFuncStmts = append(operationFuncStmts,
					generator.NewReturnStatement("nil /* TODO */", "httpRes", "nil"),
				)
			}
		} else {
			operationFuncStmts = append(operationFuncStmts,
				generator.NewReturnStatement("httpRes", "nil"),
			)
			funcSignature = funcSignature.AddReturnTypes("*http.Response", "error")
		}

		operationFunc := generator.NewFunc(
			clientStructReceiver,
			funcSignature,
			operationFuncStmts...,
		)

		clientInterface = clientInterface.AddSignatures(funcSignature)

		if summ := op.Operation.Summary; summ != "" {
			funcStmts = append(funcStmts, generator.NewComment(summ))
		}

		if desc := op.Operation.Description; desc != "" {
			funcStmts = append(funcStmts, generator.NewComment(""), generatorx.NewMultilineComment(desc))
		}

		funcStmts = append(
			funcStmts,
			operationFunc,
			generator.NewNewline(),
		)
	}

	stmts := []generator.Statement{clientInterface, clientStruct}
	stmts = append(stmts, funcStmts...)

	return stmts
}

func (c *Client) EmitReference(ctx *GeneratorContext) string {
	if c.name.PackageKey == ctx.CurrentPackage {
		return c.name.StructName
	}

	return fmt.Sprintf("%s.%s", c.name.PackageKey, c.name.StructName)
}

func (c *Client) EmitImplReference(ctx *GeneratorContext) string {
	if c.name.PackageKey == ctx.CurrentPackage {
		return c.ImplName()
	}

	return fmt.Sprintf("%s.%s", c.name.PackageKey, c.ImplName())
}
