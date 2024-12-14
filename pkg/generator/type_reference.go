package generator

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/moznion/gowrtr/generator"
)

var _ Type = &ReferenceType{}

type ReferenceType struct {
	BaseType
	Target string
}

func (o *ReferenceType) IsLightweight() bool {
	return true
}

func (o *ReferenceType) EmitDeclaration(ctx *GeneratorContext) []generator.Statement {
	output, _ := json.Marshal(o.Schema.Schema())
	return []generator.Statement{
		generator.NewComment(string(output)),
		generator.NewRawStatementf("type %s = %s", o.Names.StructName, o.EmitReference(ctx)),
	}
}

func (o *ReferenceType) EmitReference(ctx *GeneratorContext) string {
	target, err := ctx.KnownTypes.LookupReference(o.Target)
	if err == nil {
		return target.EmitReference(ctx)
	}

	log.Warn("could resolve reference", "ref", o.Target, "err", err)
	return fmt.Sprintf("ERROR /* could not resolve %s */", o.Target)
	//return "any"
}
