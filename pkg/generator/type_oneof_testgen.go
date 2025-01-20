package generator

import (
	"encoding/json"
	"fmt"
	"github.com/mittwald/api-client-go-builder/pkg/generatorx"
	"github.com/moznion/gowrtr/generator"
)

func (o *OneOfType) EmitTestCases(ctx *GeneratorContext) []generator.Statement {
	stmts := make([]generator.Statement, 0)

	suiteFunc := generator.NewAnonymousFunc(false, generator.NewAnonymousFuncSignature())
	unmarshalSuiteFunc := generator.NewAnonymousFunc(false, generator.NewAnonymousFuncSignature())

	for i, alt := range o.AlternativeTypes {
		example := alt.BuildExample(ctx, 0, 5)
		exampleJSON, _ := json.Marshal(example)

		testCaseName := fmt.Sprintf("should unmarshal into %s", o.alternativeName(i))

		testFunc := generator.NewAnonymousFunc(false, generator.NewAnonymousFuncSignature())
		testFunc = testFunc.AddStatements(
			generator.NewRawStatementf("exampleJSON := []byte(%#v)", string(exampleJSON)),
			generator.NewNewline(),
			generator.NewRawStatementf("sut := %s.%s{}", o.Names.PackageKey, o.Names.StructName),
			generator.NewRawStatementf("Expect(json.Unmarshal(exampleJSON, &sut)).To(Succeed())"),
			generator.NewRawStatementf("Expect(sut.Validate()).To(Succeed())"),
			generator.NewRawStatementf("Expect(sut.%s).NotTo(BeNil())", o.alternativeName(i)),
		)

		unmarshalSuiteFunc = unmarshalSuiteFunc.AddStatements(generatorx.NewIt(testCaseName, testFunc))
	}

	suiteFunc = suiteFunc.AddStatements(generatorx.NewWhen("unmarshaling from JSON", unmarshalSuiteFunc))

	stmts = append(stmts, generatorx.NewDescribe(o.Names.StructName, suiteFunc))

	return stmts
}
