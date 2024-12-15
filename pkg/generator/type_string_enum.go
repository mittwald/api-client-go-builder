package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
	"gopkg.in/yaml.v3"
	"strings"
)

var _ Type = &StringEnumType{}

type StringEnumType struct {
	BaseType
	Cases []string
}

func NewStringEnumTypeFromYamlNodes(baseType BaseType, nodes []*yaml.Node) *StringEnumType {
	stringCases := make([]string, len(nodes))

	for i, node := range nodes {
		stringCases[i] = node.Value
	}

	return &StringEnumType{BaseType: baseType, Cases: stringCases}
}

func (o *StringEnumType) IsLightweight() bool {
	return false
}

func (o *StringEnumType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	stmts := make([]generator.Statement, 0)

	if o.schema.Schema().Description != "" {
		stmts = append(stmts, generatorx.NewMultilineComment(o.schema.Schema().Description))
	}

	stmts = append(stmts,
		generator.NewRawStatementf("type %s string", o.Names.StructName),
		generator.NewNewline(),
	)

	caseComparisons := make([]string, len(o.Cases))

	for i, c := range o.Cases {
		caseName := c

		if caseName == "" {
			caseName = "Empty"
		}

		caseTypeName := o.Names.StructName + util.ConvertToTypename(caseName)
		stmts = append(stmts, generator.NewRawStatementf("const %s %s = %#v", caseTypeName, o.Names.StructName, c))

		caseComparisons[i] = fmt.Sprintf("e == %s", caseTypeName)
	}

	stmts = append(stmts, generator.NewFunc(
		generator.NewFuncReceiver("e", o.Names.StructName),
		generator.NewFuncSignature("Validate").
			AddReturnTypes("error"),
		generator.NewIf(strings.Join(caseComparisons, " || "), generator.NewReturnStatement("nil")),
		generator.NewReturnStatement(`fmt.Errorf("unexpected value for type %T: %s", e, e)`),
	))

	return stmts
}

func (o *StringEnumType) EmitReference(ctx *GeneratorContext) string {
	if ctx.CurrentPackage == o.Names.PackageKey {
		return o.Names.StructName
	}
	return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
}

func (o *StringEnumType) EmitValidation(ref string, _ *GeneratorContext) string {
	return fmt.Sprintf("%s.Validate()", ref)
}

func (o *StringEnumType) BuildExample() any {
	return o.Cases[0]
}
