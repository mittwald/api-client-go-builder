package generator

import (
	"encoding/json"
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
)

var _ Type = &OneOfType{}

type OneOfType struct {
	BaseType
	AlternativeTypes []Type
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

	baseSchemaName := util.LowerFirst(o.Names.StructName) + "Schema"
	for i, alt := range o.AlternativeTypes {
		schemaJson, _ := json.Marshal(alt.Schema().Schema())

		name := baseSchemaName + o.alternativeName(i)
		stmts = append(stmts, generator.NewRawStatementf("var %s = gojsonschema.NewStringLoader(%#v)", name, string(schemaJson)))
	}

	stmts = append(stmts, structType, o.emitJSONMarshalFunc(), o.emitJSONUnmarshalFunc())
	return stmts
}

func (o *OneOfType) emitJSONUnmarshalFunc() generator.Statement {
	jsonUnmarshalStmts := make([]generator.Statement, 0)
	jsonUnmarshalStmts = append(jsonUnmarshalStmts, generator.NewRawStatement("inputLoader := gojsonschema.NewBytesLoader(input)"))

	/*for i := range o.AlternativeTypes {

	}*/

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
