package generator

import "github.com/pb33f/libopenapi/datamodel/high/base"

type BaseType struct {
	Names  SchemaName
	schema *base.SchemaProxy
}

func (t *BaseType) Name() SchemaName {
	return t.Names
}

func (t *BaseType) Schema() *base.SchemaProxy {
	return t.schema
}
