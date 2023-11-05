package directive

import (
	"reflect"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/owl"
)

type (
	Directive         = owl.Directive
	Extractor         = internal.Extractor
	FormEncoder       = internal.FormEncoder
	Encoder           = internal.Encoder
	FileEncoder       = internal.FileEncoder
	BodyEncodeDecoder = internal.BodyEncodeDecoder
	DirectiveRuntime  = internal.DirectiveRuntime
)

type Decoder[T any] internal.Decoder[T]

// EncoderFunc is a function that encodes a value of type T to a string.
// It implements the Encoder interface.
type EncoderFunc[T any] func(value T) (string, error)

func (fn EncoderFunc[T]) Encode(value reflect.Value) (string, error) {
	return fn(value.Interface().(T))
}

type DecoderFunc[T any] internal.DecoderFunc[T]

func (fn DecoderFunc[T]) Decode(value string) (T, error) {
	return fn(value)
}

type FileDecoder[T any] internal.FileDecoder[T]

// RegisterEncoder registers a Encoder, which is used to encode a value of type T to a string.
func RegisterEncoder[T any](encoder EncoderFunc[T], force ...bool) {
	typ := internal.TypeOf[T]()
	internal.PanicOnError(
		internal.DefaultRegistry.RegisterEncoder(typ, encoder, force...),
	)
}

// RegisterNamedEncoder registers a Encoder with a name. The name is used to
// reference the encoder in the "encoder" directive. Ex:
//
//	func init() {
//	    RegisterNamedEncoder[time.Time]("mydate", encodeMyDate)
//	}
//
//	type ListUsersRequest struct {
//	    MemberSince time.Time `in:"query,encoder=mydate"` // use "mydate" encoder
//	}
func RegisterNamedEncoder[T any](name string, encoderFunc EncoderFunc[T], force ...bool) {
	internal.PanicOnError(
		internal.DefaultRegistry.RegisterNamedEncoder(name, encoderFunc, force...),
	)
}

// RegisterDecoder registers a FormValueDecoder. The decoder takes in a string
// value and decodes it to value of type T. Used to decode values from
// querystring, form, header. Panics on conflict types and nil decoders.
//
// NOTE: the decoder returns the decoded value as any. For best practice, the
// underlying type of the decoded value should be T, even returning *T also
// works. If the real returned value were not T or *T, the decoder will return
// an error (ErrValueTypeMismatch) while decoding.
func RegisterDecoder[T any](decoder Decoder[T], replace ...bool) {
	internal.PanicOnError(
		internal.DefaultRegistry.RegisterDecoder(
			internal.TypeOf[T](),
			internal.ToAnyDecoder[T](decoder),
			replace...,
		),
	)
}

// RegisterNamedDecoder registers a decoder by name. Panics on conflict names
// and invalid decoders. It decodes a string value to a value of type T. Use the
// "decoder" directive to override the decoder of a struct field:
//
//	func init() {
//	    httpin.RegisterNamedDecoder[time.Time]("x_time", DecoderFunc[string](decodeTimeInXFormat))
//	}
//
//	type Input struct {
//	    // httpin will use the decoder registered above, instead of the builtin decoder for time.Time.
//	    Time time.Time `in:"query:time;decoder=x_time"`
//	}
//
// Visit https://ggicci.github.io/httpin/directives/decoder for more details.
func RegisterNamedDecoder[T any](name string, decoder Decoder[T], replace ...bool) {
	internal.PanicOnError(
		internal.DefaultRegistry.RegisterNamedDecoder(
			name,
			internal.TypeOf[T](),
			internal.ToAnyDecoder[T](decoder),
			replace...,
		),
	)
}

// RegisterFileType registers a FileEncodeDecoder for type T. Which marks the type T as
// a file type. When httpin encounters a field of type T, it will treat it as a file
// upload.
//
//	func init() {
//	    RegisterFileType[MyFile](&myFileEncodeDecoder{})
//	}
func RegisterFileType[T FileEncoder](fd FileDecoder[T]) {
	internal.PanicOnError(
		internal.DefaultRegistry.RegisterFileType(
			internal.TypeOf[T](),
			internal.ToAnyFileDecoder[T](fd),
		),
	)
}

// RegisterBodyFormat registers a new data formatter for the body request, which has the
// BodyEncoderDecoder interface implemented. Panics on taken name, empty name or nil
// decoder. Pass parameter force (true) to ignore the name conflict.
//
// The BodyEncoderDecoder is used by the body directive to decode and encode the data in
// the given format (body format).
//
// It is also useful when you want to override the default registered
// BodyEncoderDecoder. For example, the default JSON decoder is borrowed from
// encoding/json. You can replace it with your own implementation, e.g.
// json-iterator/go. For example:
//
//	func init() {
//	    RegisterBodyFormat("json", &myJSONBody{}, true) // force register, replace the old one
//	    RegisterBodyFormat("yaml", &myYAMLBody{}) // register a new body format "yaml"
//	}
func RegisterBodyFormat(format string, body BodyEncodeDecoder, force ...bool) {
	internal.PanicOnError(
		internal.DefaultRegistry.RegisterBodyFormat(format, body, force...),
	)
}
