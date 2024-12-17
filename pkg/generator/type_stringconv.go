package generator

type TypeWithStringConversion interface {
	EmitToString(ref string, ctx *GeneratorContext) string
}
