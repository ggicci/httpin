package core

import (
	"errors"
	"fmt"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/owl"
)

func init() {
	RegisterDirective("form", &DirectvieForm{})
	RegisterDirective("query", &DirectiveQuery{})
	RegisterDirective("header", &DirectiveHeader{})
	RegisterDirective("body", &DirectiveBody{})
	RegisterDirective("required", &DirectiveRequired{})
	RegisterDirective("default", &DirectiveDefault{})
}

func RegisterEncoder[T any](encoderFunc EncoderFunc[T], force ...bool) {
	internal.PanicOnError(
		DefaultRegistry.RegisterEncoder(internal.TypeOf[T](), encoderFunc, force...),
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
		DefaultRegistry.RegisterNamedEncoder(name, encoderFunc, force...),
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
		DefaultRegistry.RegisterDecoder(
			internal.TypeOf[T](),
			ToAnyDecoder[T](decoder),
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
		DefaultRegistry.RegisterNamedDecoder(
			name,
			internal.TypeOf[T](),
			ToAnyDecoder[T](decoder),
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
		DefaultRegistry.RegisterFileType(
			internal.TypeOf[T](),
			toAnyFileDecoder[T](fd),
		),
	)
}

var (
	// decoderNamespace is the namespace for registering directive executors that are
	// used to decode the http request to input struct.
	decoderNamespace = owl.NewNamespace()

	// encoderNamespace is the namespace for registering directive executors that are
	// used to encode the input struct to http request.
	encoderNamespace = owl.NewNamespace()

	reservedExecutorNames = []string{"decoder", "encoder"}
	noopDirective         = &directiveNoop{}
)

type DirectiveExecutor interface {
	// Encode encodes the field of the input struct to the HTTP request.
	Encode(*DirectiveRuntime) error

	// Decode decodes the field of the input struct from the HTTP request.
	Decode(*DirectiveRuntime) error
}

func init() {
	// decoder is a special executor which does nothing, but is an indicator of
	// overriding the decoder for a specific field.
	decoderNamespace.RegisterDirectiveExecutor("decoder", asOwlDirectiveExecutor(noopDirective.Decode))
	encoderNamespace.RegisterDirectiveExecutor("encoder", asOwlDirectiveExecutor(noopDirective.Encode))
}

// RegisterDirective registers a DirectiveExecutor with the given directive name. The
// directive should be able to both extract the value from the HTTP request and build
// the HTTP request from the value. The Decode API is used to decode data from the HTTP
// request to a field of the input struct, and Encode API is used to encode the field of
// the input struct to the HTTP request.
//
// Will panic if the name were taken or given executor is nil. Pass parameter force
// (true) to ignore the name conflict.
func RegisterDirective(name string, executor DirectiveExecutor, force ...bool) {
	registerDirectiveExecutorToNamespace(decoderNamespace, name, executor, force...)
	registerDirectiveExecutorToNamespace(encoderNamespace, name, executor, force...)
}

func registerDirectiveExecutorToNamespace(ns *owl.Namespace, name string, exe DirectiveExecutor, force ...bool) {
	panicOnReservedExecutorName(name)
	if exe == nil {
		internal.PanicOnError(errors.New("nil directive executor"))
	}
	if ns == decoderNamespace {
		ns.RegisterDirectiveExecutor(name, asOwlDirectiveExecutor(exe.Decode), force...)
	} else {
		ns.RegisterDirectiveExecutor(name, asOwlDirectiveExecutor(exe.Encode), force...)
	}
}

func asOwlDirectiveExecutor(directiveFunc func(*DirectiveRuntime) error) owl.DirectiveExecutor {
	return owl.DirectiveExecutorFunc(func(dr *owl.DirectiveRuntime) error {
		return directiveFunc((*DirectiveRuntime)(dr))
	})
}

func panicOnReservedExecutorName(name string) {
	for _, reservedName := range reservedExecutorNames {
		if name == reservedName {
			internal.PanicOnError(fmt.Errorf("reserved executor name: %q", name))
		}
	}
}

// directiveNoop is a DirectiveExecutor that does nothing, "noop" stands for "no operation".
type directiveNoop struct{}

func (*directiveNoop) Encode(*DirectiveRuntime) error { return nil }
func (*directiveNoop) Decode(*DirectiveRuntime) error { return nil }
