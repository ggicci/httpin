package httpin

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func encodeCustomBool(value bool) (string, error) {
	if value {
		return "yes", nil
	}
	return "no", nil
}

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

	enc := toPointerEncoder{EncoderFunc[bool](encodeCustomBool)}

	actual, err = enc.Encode(reflect.ValueOf(asPointerValue[bool](true)))
	assert.NoError(err)
	assert.Equal("yes", actual)

	actual, err = enc.Encode(reflect.ValueOf(asPointerValue[bool](false)))
	assert.NoError(err)
	assert.Equal("no", actual)
}

func removeTypeEncoder[T any]() {
	delete(customEncoders, typeOf[T]())
	delete(customEncoders, typeOf[*T]())
}

func removeNamedEncoder(name string) {
	delete(namedEncoders, name)
}

func asPointerValue[T any](v T) *T {
	return &v
}
