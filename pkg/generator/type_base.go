package generator

import "github.com/pb33f/libopenapi/datamodel/high/base"

type BaseType struct {
	Names  SchemaName
	Schema *base.SchemaProxy
}

func (t *BaseType) Name() SchemaName {
	return t.Names
}
