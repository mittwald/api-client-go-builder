package generator

// UnpackableType describes a data type that is a wrapper around another
// type. An example for this is the OptionalType wrapper.
type UnpackableType interface {
	Unpack() SchemaType
}
