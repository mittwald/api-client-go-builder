package generator

import (
	"encoding/json"
	"github.com/moznion/gowrtr/generator"
)

var _ SchemaType = &UnknownType{}

// UnknownType is a placeholder for any type that could not be automatically
// generated from a JSON schema. It is used as fallback type in various places,
// but should ideally never be used.
type UnknownType struct {
	BaseType
}

func (o *UnknownType) IsLightweight() bool {
	return true
}

func (o *UnknownType) EmitDeclaration(*GeneratorContext) []generator.Statement {
	output, _ := json.Marshal(o.schema.Schema())
	return []generator.Statement{
		generator.NewCommentf("TODO: This schema could not be automatically generated"),
		generator.NewComment(string(output)),
		generator.NewRawStatementf("type %s = any", o.Names.StructName),
	}
}

func (o *UnknownType) EmitReference(*GeneratorContext) string {
	//return fmt.Sprintf("%s.%s", o.Names.PackageKey, o.Names.StructName)
	return "any"
}

func (o *UnknownType) BuildExample(*GeneratorContext, int, int) any {
	return nil
}

func (o *UnknownType) IsPointerType() bool {
	// Even though `any` can be a pointer type, we treat it as non-pointer type here;
	// since `nil` is a valid value, we still need to use a pointer to `any` when
	// the field is optional.
	return false
}
