package generator

import (
	"encoding/json"
	"fmt"
	"github.com/moznion/gowrtr/generator"
)

func (o *OneOfType) EmitTestCases(ctx *GeneratorContext) []generator.Statement {
	stmts := make([]generator.Statement, 0)

	for i, alt := range o.AlternativeTypes {
		example := alt.BuildExample(ctx)
		exampleJSON, _ := json.Marshal(example)

		testCaseName := fmt.Sprintf("Test%sCanUnmarshal%s", o.Names.StructName, o.alternativeName(i))

		testFunc := generator.NewFunc(
			nil,
			generator.NewFuncSignature(testCaseName).
				AddParameters(generator.NewFuncParameter("t", "*testing.T")),
			generator.NewRawStatementf("exampleJSON := []byte(%#v)", string(exampleJSON)),
			generator.NewNewline(),
			generator.NewRawStatementf("sut := %s{}", o.Names.StructName),
			generator.NewIf("err := json.Unmarshal(exampleJSON, &sut); err != nil", generator.NewRawStatement("t.Fatalf(\"could not unmarshal: %s\", err.Error())")),
			generator.NewIf("err := sut.Validate(); err != nil", generator.NewRawStatement("t.Fatalf(\"could not unmarshal to a valid struct: %s\", err.Error())")),
			generator.NewIf(fmt.Sprintf("sut.%s == nil", o.alternativeName(i)), generator.NewRawStatement("t.Fatal(\"expected alternative was nil\")")),
		)

		for j := range o.AlternativeTypes {
			if j != i {
				testFunc = testFunc.AddStatements(
					generator.NewIf(fmt.Sprintf("sut.%s != nil", o.alternativeName(j)), generator.NewRawStatement("t.Fatal(\"unexpected alternative was not nil\")")),
				)
			}
		}

		stmts = append(stmts, testFunc, generator.NewNewline())
	}

	return stmts
}
