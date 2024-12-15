package generator

type GeneratorContext struct {
	CurrentPackage        string
	KnownTypes            *TypeStore
	WithDebuggingComments bool
}
