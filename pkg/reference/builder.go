package reference

import (
	"fmt"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"strings"
)

type ReferenceLinkBuilder func(op *v3.Operation) (string, bool)

func NewMittwaldReferenceLinkBuilder(apiVersion string) ReferenceLinkBuilder {
	return func(op *v3.Operation) (string, bool) {
		return fmt.Sprintf(
			"https://developer.mittwald.de/docs/%s/reference/%s/%s",
			apiVersion,
			strings.ToLower(op.Tags[0]),
			op.OperationId,
		), true
	}
}
