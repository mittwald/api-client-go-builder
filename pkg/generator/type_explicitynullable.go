package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/moznion/gowrtr/generator"
)

var _ SchemaType = &ExplicitlyNullableType{}

type ExplicitlyNullableType struct {
	BaseType

	InnerType SchemaType
}

func (o *ExplicitlyNullableType) IsLightweight() bool {
	return false
	//return true
}

func (o *ExplicitlyNullableType) BuildSubtypes(opts GeneratorOpts, store *TypeStore) error {
	if s, ok := o.InnerType.(TypeWithSubtypes); ok {
		return s.BuildSubtypes(opts, store)
	}
	return nil
}

func (o *ExplicitlyNullableType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	stmts := make([]generator.Statement, 0)

	structType := generator.NewStruct(o.Names.StructName)
	structType = structType.AddField("Value", "*"+o.InnerType.EmitReference(ctx))

	stmts = append(stmts,
		generatorx.NewWrappingCommentf("%s is a wrapper around %s, which allows you to define an explicit NULL value. This is useful for PATCH routes, in which an explicit NULL value may have a different semantic that a missing value.", o.Names.StructName, o.InnerType.EmitReference(ctx)),
		structType,
		o.emitJSONMarshalFunc(),
		generator.NewNewline(),
		o.emitJSONUnmarshalFunc(ctx),
		generator.NewNewline(),
		o.emitValidateFunc(ctx),
	)
	return stmts
}

func (o *ExplicitlyNullableType) emitJSONMarshalFunc() generator.Statement {
	jsonMarshalStmts := make([]generator.Statement, 0)
	jsonMarshalStmts = append(
		jsonMarshalStmts,
		generator.NewIf(
			"a.Value != nil",
			generator.NewReturnStatement("json.Marshal(a.Value)"),
		),
	)
	jsonMarshalStmts = append(jsonMarshalStmts, generator.NewReturnStatement(`[]byte("null"), nil`))

	return generator.NewFunc(
		generator.NewFuncReceiver("a", fmt.Sprintf("*%s", o.Names.StructName)),
		generator.NewFuncSignature("MarshalJSON").
			AddReturnTypes("[]byte", "error"),
		jsonMarshalStmts...,
	)
}

func (o *ExplicitlyNullableType) emitValidateFunc(ctx *GeneratorContext) generator.Statement {
	validateStmts := make([]generator.Statement, 0)

	if iv, ok := o.InnerType.(TypeWithValidation); ok {
		validateStmts = append(
			validateStmts,
			generator.NewIf(
				"a.Value != nil",
				generator.NewReturnStatement(iv.EmitValidation("a.Value", ctx)),
			),
		)
	}

	validateStmts = append(validateStmts, generator.NewReturnStatement(`nil`))

	return generator.NewFunc(
		generator.NewFuncReceiver("a", fmt.Sprintf("*%s", o.Names.StructName)),
		generator.NewFuncSignature("Validate").
			AddReturnTypes("error"),
		validateStmts...,
	)
}

func (o *ExplicitlyNullableType) emitJSONUnmarshalFunc(ctx *GeneratorContext) generator.Statement {
	jsonUnmarshalStmts := make([]generator.Statement, 0)
	jsonUnmarshalStmts = append(jsonUnmarshalStmts,
		generator.NewIf("string(input) == \"null\"",
			generator.NewRawStatement("a.Value = nil"),
			generator.NewReturnStatement("nil"),
		),
	)

	jsonUnmarshalStmts = append(jsonUnmarshalStmts,
		generator.NewNewline(),
		generator.NewRawStatementf("a.Value = new(%s)", o.InnerType.EmitReference(ctx)),
		generator.NewReturnStatement("json.Unmarshal(input, a.Value)"),
	)

	return generator.NewFunc(
		generator.NewFuncReceiver("a", fmt.Sprintf("*%s", o.Names.StructName)),
		generator.NewFuncSignature("UnmarshalJSON").
			AddParameters(generator.NewFuncParameter("input", "[]byte")).
			AddReturnTypes("error"),
		jsonUnmarshalStmts...,
	)
}

func (o *ExplicitlyNullableType) EmitReference(ctx *GeneratorContext) string {
	if ctx.CurrentPackage == o.Names.PackageKey {
		return o.Names.StructName
	}
	return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
}

func (o *ExplicitlyNullableType) EmitValidation(ref string, ctx *GeneratorContext) string {
	if v, ok := o.InnerType.(TypeWithValidation); ok {
		return fmt.Sprintf("func () error {\nif %s == nil {\nreturn nil\n}\nreturn %s\n}()", ref, v.EmitValidation(ref, ctx))
	}
	return "nil"
}

func (o *ExplicitlyNullableType) BuildExample(ctx *GeneratorContext, level, maxLevel int) any {
	if level == maxLevel {
		return nil
	}

	return o.InnerType.BuildExample(ctx, level+1, maxLevel)
}

func (o *ExplicitlyNullableType) EmitToString(ref string, ctx *GeneratorContext) string {
	if ts, ok := o.InnerType.(TypeWithStringConversion); ok {
		return ts.EmitToString("*"+ref, ctx)
	}

	// if they want compile errors, give them compile errors!
	return "invalid-no-string-conversion"
}

func (o *ExplicitlyNullableType) Unpack() SchemaType {
	return o.InnerType
}
