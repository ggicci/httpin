package directive

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

	enc := internal.ToPointerEncoder{Encoder: EncoderFunc[bool](encodeCustomBool)}

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

func removeTypeEncoder[T any]() {
	internal.DefaultRegistry.RemoveEncoder(internal.TypeOf[T]())
	internal.DefaultRegistry.RemoveEncoder(internal.TypeOf[*T]())
}

func removeNamedEncoder(name string) {
	internal.DefaultRegistry.RemoveNamedEncoder(name)
}

func removeTypeDecoder[T any]() {
	internal.DefaultRegistry.RemoveDecoder(internal.TypeOf[T]())
	internal.DefaultRegistry.RemoveDecoder(internal.TypeOf[[]T]())
	internal.DefaultRegistry.RemoveDecoder(internal.TypeOf[patch.Field[T]]())
	internal.DefaultRegistry.RemoveDecoder(internal.TypeOf[patch.Field[[]T]]())
}

func removeNamedDecoder(name string) {
	internal.DefaultRegistry.RemoveNamedDecoder(name)
}
