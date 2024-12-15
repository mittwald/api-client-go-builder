package generator

import (
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/util"
	"github.com/moznion/gowrtr/generator"
	"gopkg.in/yaml.v3"
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
	stmts := []generator.Statement{
		generator.NewRawStatementf("type %s = string", o.Names.StructName),
		generator.NewNewline(),
	}

	for _, c := range o.Cases {
		caseName := c

		if caseName == "" {
			caseName = "Empty"
		}

		caseTypeName := o.Names.StructName + util.ConvertToTypename(caseName)
		stmts = append(stmts, generator.NewRawStatementf("const %s %s = %#v", caseTypeName, o.Names.StructName, c))
	}

	return stmts
}

func (o *StringEnumType) EmitReference(ctx *GeneratorContext) string {
	if ctx.CurrentPackage == o.Names.PackageKey {
		return o.Names.StructName
	}
	return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
}
