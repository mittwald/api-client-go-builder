package generator

import (
	"github.com/mittwald/api-client-go-builder/pkg/reference"
)

type GeneratorContext struct {
	CurrentPackage        string
	KnownTypes            *TypeStore
	WithDebuggingComments bool
	BuildReferenceLink    reference.ReferenceLinkBuilder
}
