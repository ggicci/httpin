package core

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

func encodeCustomBool(value bool) (string, error) {
	if value {
		return "yes", nil
	}
	return "no", nil
}

// decodeCustomBool additionally parses "yes/no".
func decodeCustomBool(value string) (bool, error) {
	sdata := strings.ToLower(value)
	if sdata == "yes" {
		return true, nil
	}
	if sdata == "no" {
		return false, nil
	}
	return strconv.ParseBool(sdata)
}

var myBoolDecoder = DecoderFunc[bool](decodeCustomBool)

func TestRegisterTypeEncoder(t *testing.T) {
	assert.Panics(t, func() {
		RegisterEncoder[int](nil) // fail to register nil encoder
	})

	assert.NotPanics(t, func() {
		RegisterEncoder[bool](encodeCustomBool)
	})
	assert.Panics(t, func() {
		RegisterEncoder[bool](encodeCustomBool) // fail to register duplicate encoder
	})

	removeTypeEncoder[bool]()
}

func TestRegisterTypeEncoder_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterEncoder[bool](encodeCustomBool, true)
	})

	assert.NotPanics(t, func() {
		RegisterEncoder[bool](encodeCustomBool, true)
	})

	removeTypeEncoder[bool]()
}

func TestRegisterNamedEncoder(t *testing.T) {
	assert.Panics(t, func() {
		RegisterNamedEncoder[bool]("myBool", nil) // fail to register nil encoder
	})

	assert.NotPanics(t, func() {
		RegisterNamedEncoder[bool]("myBool", encodeCustomBool)
	})
	assert.Panics(t, func() {
		RegisterNamedEncoder[bool]("myBool", encodeCustomBool) // fail to register duplicate encoder
	})

	removeNamedEncoder("myBool")
}

func TestRegisterNamedEncoder_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterNamedEncoder[bool]("myBool", encodeCustomBool, true)
	})

	assert.NotPanics(t, func() {
		RegisterNamedEncoder[bool]("myBool", encodeCustomBool, true)
	})

	removeNamedEncoder("myBool")
}

func Test_scaler2pointerEncoder(t *testing.T) {
	assert := assert.New(t)

	actual, err := encodeCustomBool(true)
	assert.NoError(err)
	assert.Equal("yes", actual)

	actual, err = encodeCustomBool(false)
	assert.NoError(err)
	assert.Equal("no", actual)

	enc := ToPointerEncoder{Encoder: EncoderFunc[bool](encodeCustomBool)}

	actual, err = enc.Encode(reflect.ValueOf(internal.Pointerize[bool](true)))
	assert.NoError(err)
	assert.Equal("yes", actual)

	actual, err = enc.Encode(reflect.ValueOf(internal.Pointerize[bool](false)))
	assert.NoError(err)
	assert.Equal("no", actual)
}

func TestRegisterValueTypeDecoder(t *testing.T) {
	assert.Panics(t, func() { RegisterDecoder[bool](nil) }) // fail on nil decoder

	assert.NotPanics(t, func() {
		RegisterDecoder[bool](myBoolDecoder)
	})
	assert.Panics(t, func() {
		// Fail on duplicate registeration on the same type.
		RegisterDecoder[bool](myBoolDecoder)
	})
	removeTypeDecoder[bool]() // remove the custom decoder
}

func TestRegisterValueTypeDecoder_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterDecoder[bool](myBoolDecoder, true)
	})

	assert.NotPanics(t, func() {
		RegisterDecoder[bool](myBoolDecoder, true)
	})

	removeTypeDecoder[bool]() // remove the custom decoder
}

func TestRegisterNamedDecoder(t *testing.T) {
	assert.Panics(t, func() { RegisterNamedDecoder[bool]("myBool", nil) }) // fail on nil decoder

	// Register duplicate decoder should fail.
	assert.NotPanics(t, func() {
		RegisterNamedDecoder[bool]("mybool", myBoolDecoder)
	})
	assert.Panics(t, func() {
		// Fail on duplicate registeration on the same name.
		RegisterNamedDecoder[bool]("mybool", myBoolDecoder)
	})

	removeNamedDecoder("mybool") // remove the custom decoder
}

func TestRegisterNamedDecoder_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterNamedDecoder[bool]("mybool", myBoolDecoder, true)
	})

	assert.NotPanics(t, func() {
		RegisterNamedDecoder[bool]("mybool", myBoolDecoder, true)
	})

	removeNamedDecoder("mybool") // remove the custom decoder
}

func removeType[T any]() {
	delete(customStringableAdaptors, internal.TypeOf[T]())
}

func removeTypeEncoder[T any]() {
	defaultRegistry.RemoveEncoder(internal.TypeOf[T]())
	defaultRegistry.RemoveEncoder(internal.TypeOf[*T]())
}

func removeNamedEncoder(name string) {
	defaultRegistry.RemoveNamedEncoder(name)
}

func removeTypeDecoder[T any]() {
	defaultRegistry.RemoveDecoder(internal.TypeOf[T]())
	defaultRegistry.RemoveDecoder(internal.TypeOf[[]T]())
	defaultRegistry.RemoveDecoder(internal.TypeOf[patch.Field[T]]())
	defaultRegistry.RemoveDecoder(internal.TypeOf[patch.Field[[]T]]())
}

func removeNamedDecoder(name string) {
	defaultRegistry.RemoveNamedDecoder(name)
}

func TestRegisterDirectiveExecutor(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterDirective("noop_TestRegisterDirectiveExecutor", noopDirective)
	})

	assert.Panics(t, func() {
		RegisterDirective("noop_TestRegisterDirectiveExecutor", noopDirective)
	}, "should panic on duplicate name")

	assert.Panics(t, func() {
		RegisterDirective("nil_TestRegisterDirectiveExecutor", nil)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		RegisterDirective("decoder", noopDirective)
	}, "should panic on reserved name")
}

func TestRegisterDirectiveExecutor_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterDirective("noop_TestRegisterDirectiveExecutor_forceReplace", noopDirective, true)
	})

	assert.NotPanics(t, func() {
		RegisterDirective("noop_TestRegisterDirectiveExecutor_forceReplace", noopDirective, true)
	}, "should not panic on duplicate name")

	assert.Panics(t, func() {
		RegisterDirective("nil_TestRegisterDirectiveExecutor_forceReplace", nil, true)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		RegisterDirective("decoder", noopDirective, true)
	}, "should panic on reserved name")
}
