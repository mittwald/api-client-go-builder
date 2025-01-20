package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
)

var _ SchemaType = &OneOfType{}

type OneOfType struct {
	BaseType
	AlternativeTypes []SchemaType
}

func (o *OneOfType) BuildSubtypes(store *TypeStore) error {
	for i, alt := range o.AlternativeTypes {
		subTypeName := o.Names.ForSubtype(o.alternativeName(i))
		store.AddSubtype(subTypeName, alt)
	}

	return nil
}

func (o *OneOfType) alternativeName(idx int) string {
	alternativeName := o.AlternativeTypes[idx].Name().StructName
	return fmt.Sprintf("Alternative%s", alternativeName)
	//return fmt.Sprintf("Alternative%d", idx+1)
}

func (o *OneOfType) IsLightweight() bool {
	return false
}

func (o *OneOfType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	stmts := make([]generator.Statement, 0)

	structType := generator.NewStruct(o.Names.StructName)
	for i, alt := range o.AlternativeTypes {
		name := o.alternativeName(i)
		structType = structType.AddField(name, "*"+alt.EmitReference(ctx))
	}

	stmts = append(stmts,
		structType,
		o.emitJSONMarshalFunc(),
		generator.NewNewline(),
		o.emitJSONUnmarshalFunc(ctx),
		generator.NewNewline(),
		o.emitValidateFunc(ctx),
	)
	return stmts
}

func (o *OneOfType) emitValidateFunc(ctx *GeneratorContext) generator.Statement {
	validateStmts := make([]generator.Statement, 0)

	for i, alt := range o.AlternativeTypes {
		ref := fmt.Sprintf("a.%s", o.alternativeName(i))
		if v, ok := alt.(TypeWithValidation); ok {
			altValidation := v.EmitValidation(ref, ctx)
			validateStmts = append(validateStmts,
				generator.NewIf(ref+" != nil", generator.NewReturnStatement(altValidation)),
			)
		} else {
			validateStmts = append(validateStmts,
				generator.NewCommentf(" The %s subtype does not implement validation, so we consider being non-nil as valid", o.alternativeName(i)),
				generator.NewIf(ref+" != nil", generator.NewReturnStatement("nil")),
			)
		}
	}

	validateStmts = append(validateStmts, generator.NewReturnStatement("errors.New(\"no alternative set\")"))

	return generator.NewFunc(
		generator.NewFuncReceiver("a", fmt.Sprintf("*%s", o.Names.StructName)),
		generator.NewFuncSignature("Validate").
			AddReturnTypes("error"),
		validateStmts...,
	)
}

func (o *OneOfType) emitJSONUnmarshalFunc(ctx *GeneratorContext) generator.Statement {
	jsonUnmarshalStmts := make([]generator.Statement, 0)

	jsonUnmarshalStmts = append(jsonUnmarshalStmts,
		generator.NewRawStatement("reader := bytes.NewReader(input)"),
		generator.NewRawStatement("decodedAtLeastOnce := false"),
		generator.NewRawStatement("dec := json.NewDecoder(reader)"),
		generator.NewRawStatement("dec.DisallowUnknownFields()"),
		generator.NewNewline(),
	)

	for i, alt := range o.AlternativeTypes {
		name := o.alternativeName(i)
		localName := util.LowerFirst(name)

		unmarshalCondition := generator.NewIf(
			fmt.Sprintf("err := dec.Decode(&%s); err == nil", localName),
		)

		if v, ok := alt.(TypeWithValidation); ok {
			validation := v.EmitValidation(localName, ctx)
			if validation != "nil" {
				unmarshalCondition = unmarshalCondition.AddStatements(
					generator.NewCommentf("subtype: %T", alt),
					generator.NewIf(fmt.Sprintf("vErr := %s; vErr == nil", validation),
						generator.NewRawStatementf("a.%s = &%s", name, localName),
						generator.NewRawStatement("decodedAtLeastOnce = true"),
					),
				)
			} else {
				unmarshalCondition = unmarshalCondition.AddStatements(
					generator.NewCommentf("subtype: %T", alt),
					generator.NewRawStatementf("a.%s = &%s", name, localName),
					generator.NewRawStatement("decodedAtLeastOnce = true"),
				)
			}
		} else {
			unmarshalCondition = unmarshalCondition.AddStatements(
				generator.NewRawStatementf("a.%s = &%s", name, localName),
				generator.NewRawStatement("decodedAtLeastOnce = true"),
			)
		}

		jsonUnmarshalStmts = append(jsonUnmarshalStmts,
			generator.NewRawStatement("reader.Reset(input)"),
			generator.NewRawStatementf("var %s %s", localName, alt.EmitReference(ctx)),
			unmarshalCondition,
			generator.NewNewline(),
		)
	}

	jsonUnmarshalStmts = append(jsonUnmarshalStmts,
		generator.NewIf("!decodedAtLeastOnce",
			generator.NewReturnStatement("fmt.Errorf(\"could not unmarshal into any alternative for type %T\", a)")),
		generator.NewReturnStatement("nil"),
	)

	return generator.NewFunc(
		generator.NewFuncReceiver("a", fmt.Sprintf("*%s", o.Names.StructName)),
		generator.NewFuncSignature("UnmarshalJSON").
			AddParameters(generator.NewFuncParameter("input", "[]byte")).
			AddReturnTypes("error"),
		jsonUnmarshalStmts...,
	)
}

func (o *OneOfType) emitJSONMarshalFunc() generator.Statement {
	jsonMarshalStmts := make([]generator.Statement, 0)
	for i := range o.AlternativeTypes {
		jsonMarshalStmts = append(
			jsonMarshalStmts,
			generator.NewIf(
				fmt.Sprintf("a.%s != nil", o.alternativeName(i)),
				generator.NewReturnStatement(fmt.Sprintf("json.Marshal(a.%s)", o.alternativeName(i))),
			),
		)
	}
	jsonMarshalStmts = append(jsonMarshalStmts, generator.NewReturnStatement(`[]byte("null"), nil`))

	return generator.NewFunc(
		generator.NewFuncReceiver("a", fmt.Sprintf("*%s", o.Names.StructName)),
		generator.NewFuncSignature("MarshalJSON").
			AddReturnTypes("[]byte", "error"),
		jsonMarshalStmts...,
	)
}

func (o *OneOfType) EmitReference(ctx *GeneratorContext) string {
	if ctx.CurrentPackage == o.Names.PackageKey {
		return o.Names.StructName
	}

	return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
}

func (o *OneOfType) EmitValidation(ref string, ctx *GeneratorContext) string {
	return ref + ".Validate()"
}

func (o *OneOfType) BuildExample(ctx *GeneratorContext, level, maxLevel int) any {
	return o.AlternativeTypes[0].BuildExample(ctx, level, maxLevel)
}
