package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"regexp"
	"strings"
)

type ClientOperationRequest struct {
	name      SchemaName
	operation *OperationWithMeta

	bodyType   Type
	bodyFormat string

	pathParams  *orderedmap.OrderedMap[string, SchemaType]
	queryParams *orderedmap.OrderedMap[string, SchemaType]
}

func (c *ClientOperationRequest) Name() SchemaName {
	return c.name
}

func (c *ClientOperationRequest) BuildSubtypes(_ GeneratorOpts, store *TypeStore) error {
	c.pathParams = orderedmap.New[string, SchemaType]()
	c.queryParams = orderedmap.New[string, SchemaType]()

	for _, param := range c.operation.Operation.Parameters {
		tmp := c.name.ForSubtype(param.In)
		//paramFieldName := util.ConvertToTypename(param.Name)
		paramName := tmp.ForSubtype(param.Name)
		paramType, err := BuildTypeFromSchema(paramName, param.Schema, store)
		if err != nil {
			return err
		}

		if param.Required == nil || *param.Required == false {
			paramType = &OptionalType{InnerType: paramType}
		}

		if param.In == "path" {
			c.pathParams.Set(param.Name, paramType)
		} else if param.In == "query" {
			c.queryParams.Set(param.Name, paramType)
		}

		store.AddSubtype(paramName, paramType)
	}

	if c.operation.Operation.RequestBody != nil {
		if jsonBody, ok := c.operation.Operation.RequestBody.Content.Get("application/json"); ok {
			bodyName := c.name.ForSubtype("Body")
			bodyType, err := BuildTypeFromSchema(bodyName, jsonBody.Schema, store)
			if err != nil {
				return err
			}

			schema := jsonBody.Schema.Schema()
			if schema.Description == "" {
				schema.Description = fmt.Sprintf("%s models the JSON body of a '%s' request", bodyName.StructName, c.operation.Operation.OperationId)
			}

			c.bodyType = bodyType
			c.bodyFormat = "json"

			store.AddSubtype(bodyName, bodyType)
		}
	}

	return nil
}

func (c *ClientOperationRequest) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	str := generator.NewStruct(c.name.StructName)

	if c.bodyType != nil {
		str = str.AddField("Body", c.bodyType.EmitReference(ctx))
	}

	for name, param := range c.pathParams.FromOldest() {
		fieldName := util.ConvertToTypename(name)

		if d := param.Schema().Schema().Description; d != "" {
			str = generatorx.AddFieldComment(str, d)
		}
		str = str.AddField(fieldName, param.EmitReference(ctx))
	}
	for name, param := range c.queryParams.FromOldest() {
		fieldName := util.ConvertToTypename(name)

		if d := param.Schema().Schema().Description; d != "" {
			str = generatorx.AddFieldComment(str, d)
		}
		str = str.AddField(fieldName, param.EmitReference(ctx))
	}

	reqFunc := c.buildNewRequestFunction()
	bodyFunc := c.buildBodyFunction()
	urlFunc := c.buildURLFunction(ctx)
	queryFunc := c.buildQueryFunction(ctx)

	strComment := c.buildDocumentation(ctx)

	return []generator.Statement{
		strComment,
		str,
		generatorx.NewWrappingComment("BuildRequest builds an *http.Request instance from this request that may be used with any regular *http.Client instance."),
		reqFunc, generator.NewNewline(),
		bodyFunc, generator.NewNewline(),
		urlFunc, generator.NewNewline(),
		queryFunc,
	}
}

func (c *ClientOperationRequest) buildDocumentation(ctx *GeneratorContext) generator.Statement {
	strComment := generatorx.NewWrappingCommentf("%s models a request for the '%s' operation.", c.name.StructName, c.operation.Operation.OperationId)
	if ctx.BuildReferenceLink != nil {
		if link, ok := ctx.BuildReferenceLink(c.operation.Operation); ok {
			strComment.Writef(" See [1] for more information.")
			defer strComment.Writef("\n\n[1]: %s", link)
		}
	}

	if s := c.operation.Operation.Summary; s != "" {
		strComment.Writef("\n\n%s", s)
	}

	if d := c.operation.Operation.Description; d != "" {
		strComment.Writef("\n\n%s", d)
	}

	return strComment
}

func (c *ClientOperationRequest) buildNewRequestFunction() generator.Statement {
	receiver := generator.NewFuncReceiver("r", "*"+c.name.StructName)
	signature := generator.NewFuncSignature("BuildRequest").AddReturnTypes("*http.Request", "error")

	method := fmt.Sprintf("http.Method%s", util.UpperFirst(c.operation.Method))

	errorReturn := generator.NewReturnStatement("nil", "err")
	funcStmts := []generator.Statement{
		generator.NewRawStatement("body, contentType, err := r.body()"),
		generator.NewIf("err != nil", errorReturn),
		generator.NewNewline(),
		//generator.NewReturnStatement(fmt.Sprintf("http.NewRequest(%s, r.url(), body)", method)),
		generator.NewRawStatementf("req, err := http.NewRequest(%s, r.url(), body)", method),
		generator.NewIf("err != nil", errorReturn),
		generator.NewRawStatement("req.Header.Set(\"Content-Type\", contentType)"),
		generator.NewReturnStatement("req", "nil"),
	}

	return generator.NewFunc(receiver, signature, funcStmts...)
}

func (c *ClientOperationRequest) buildBodyFunction() generator.Statement {
	receiver := generator.NewFuncReceiver("r", "*"+c.name.StructName)
	signature := generator.NewFuncSignature("body").AddReturnTypes("io.Reader", "string", "error")

	if c.bodyFormat == "json" {
		return generator.NewFunc(receiver, signature,
			generator.NewRawStatement("out, err := json.Marshal(&r.Body)"),
			generator.NewIf("err != nil", generator.NewReturnStatement("nil", `""`, `fmt.Errorf("error while marshalling JSON: %w", err)`)),
			generator.NewReturnStatement("bytes.NewReader(out)", `"application/json"`, "nil"),
		)
	}

	return generator.NewFunc(receiver, signature, generator.NewReturnStatement("nil", `""`, "nil"))
}

func (c *ClientOperationRequest) buildQueryFunction(ctx *GeneratorContext) generator.Statement {
	if c.queryParams.Len() == 0 {
		return generator.NewFunc(
			generator.NewFuncReceiver("r", "*"+c.name.StructName),
			generator.NewFuncSignature("query").AddReturnTypes("url.Values"),
			generator.NewReturnStatement("nil"),
		)
	}

	stmts := []generator.Statement{
		generator.NewRawStatement("q := make(url.Values)"),
	}

	for name, param := range c.queryParams.FromOldest() {
		fieldName := util.ConvertToTypename(name)

		unpacked := param
		if unpackable, ok := param.(UnpackableType); ok {
			unpacked = unpackable.Unpack()
		}

		// NOTE: Some of the type assertions below are asserted against
		// "unpacked", and some against "param". This is done on purpose.

		if paramArray, isArray := unpacked.(*ArrayType); isArray {
			if paramArrayItemString, isStringConvertible := paramArray.ItemType.(TypeWithStringConversion); isStringConvertible {
				stmt := generator.NewFor("_, val := range r."+fieldName,
					generator.NewRawStatementf("q.Add(%#v, %s)", name, paramArrayItemString.EmitToString("val", ctx)),
				)
				stmts = append(stmts, stmt)
			}
		} else if t, ok := param.(TypeWithStringConversion); ok {
			var stmt generator.Statement = generator.NewRawStatementf("q.Set(%#v, %s)", name, t.EmitToString("r."+fieldName, ctx))
			if _, isOptional := param.(*OptionalType); isOptional {
				stmt = generator.NewIf(fmt.Sprintf("r.%s != nil", fieldName), stmt)
			}

			stmts = append(stmts, stmt)
		}
	}

	stmts = append(stmts, generator.NewReturnStatement("q"))

	queryFunc := generator.NewFunc(
		generator.NewFuncReceiver("r", "*"+c.name.StructName),
		generator.NewFuncSignature("query").AddReturnTypes("url.Values"),
		stmts...,
	)

	return queryFunc
}

func (c *ClientOperationRequest) buildURLFunction(ctx *GeneratorContext) generator.Statement {
	paramPattern := regexp.MustCompile("{([a-zA-Z0-9]+)}")

	builtUrlParams := []string{}
	builtUrl := paramPattern.ReplaceAllStringFunc(c.operation.Path, func(match string) string {
		name := strings.Trim(match, "{}")
		paramName := util.ConvertToTypename(name)
		param := c.pathParams.Value(name)

		if ts, ok := param.(TypeWithStringConversion); ok {
			builtUrlParams = append(builtUrlParams, "url.PathEscape("+ts.EmitToString("r."+paramName, ctx)+")")
		}

		return "%s"
	})

	builtUrlStmt := fmt.Sprintf("%#v", builtUrl)
	if c.pathParams.Len() > 0 {
		builtUrlStmt = fmt.Sprintf("fmt.Sprintf(%#v, %s)", builtUrl, strings.Join(builtUrlParams, ", "))
	}

	urlFunc := generator.NewFunc(
		generator.NewFuncReceiver("r", "*"+c.name.StructName),
		generator.NewFuncSignature("url").AddReturnTypes("string"),
		generator.NewReturnStatement(builtUrlStmt),
	)

	return urlFunc
}

func (c *ClientOperationRequest) EmitReference(ctx *GeneratorContext) string {
	if c.name.PackageKey == ctx.CurrentPackage {
		return c.name.StructName
	}

	return fmt.Sprintf("%s.%s", c.name.PackageKey, c.name.StructName)
}
