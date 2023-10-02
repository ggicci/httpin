package httpin

import (
	"errors"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

// decodeCustomBool additionally parses "yes/no".
func decodeCustomBool(value string) (interface{}, error) {
	sdata := strings.ToLower(value)
	if sdata == "yes" {
		return true, nil
	}
	if sdata == "no" {
		return false, nil
	}
	return strconv.ParseBool(sdata)
}

var myBoolDecoder = DecoderFunc[string](decodeCustomBool)

type Place struct {
	Country string
	City    string
}

// decodePlace parses "country.city", e.g. "Canada.Toronto".
// It returns a Place.
func decodePlace(value string) (interface{}, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid place")
	}
	return Place{Country: parts[0], City: parts[1]}, nil
}

// decodePlacePointer parses "country.city", e.g. "Canada.Toronto".
// It returns *Place.
func decodePlacePointer(value string) (interface{}, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid place")
	}
	return &Place{Country: parts[0], City: parts[1]}, nil
}

var myPlaceDecoder = DecoderFunc[string](decodePlace)
var myPlacePointerDecoder = DecoderFunc[string](decodePlacePointer)

type BadFile struct{}

// decodeBadFile always returns an error, to simulate the case that we cannot
// decode the file properly.
type badFileDecoder struct{}

var errBadFile = errors.New("bad file")

func (badFileDecoder) Decode(*multipart.FileHeader) (interface{}, error) {
	return nil, errBadFile
}

func TestRegisterValueTypeDecoder(t *testing.T) {
	assert.Panics(t, func() { RegisterValueTypeDecoder[bool](nil) }) // fail on nil decoder

	assert.NotPanics(t, func() {
		RegisterValueTypeDecoder[bool](myBoolDecoder)
	})
	assert.Panics(t, func() {
		// Fail on duplicate registeration on the same type.
		RegisterValueTypeDecoder[bool](myBoolDecoder)
	})
	removeTypeDecoder[bool]() // remove the custom decoder
}

func TestRegisterValueTypeDecoder_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterValueTypeDecoder[bool](myBoolDecoder, true)
	})

	assert.NotPanics(t, func() {
		RegisterValueTypeDecoder[bool](myBoolDecoder, true)
	})

	removeTypeDecoder[bool]() // remove the custom decoder
}

func TestRegisterFileTypeDecoder(t *testing.T) {
	assert.Panics(t, func() { RegisterFileTypeDecoder[BadFile](nil) }) // fail on nil decoder

	assert.NotPanics(t, func() {
		RegisterFileTypeDecoder[BadFile](badFileDecoder{})
	})
	assert.Panics(t, func() {
		// Fail on duplicate registeration on the same type.
		RegisterFileTypeDecoder[BadFile](badFileDecoder{})
	})

	removeTypeDecoder[BadFile]() // remove the custom decoder
}

func TestRegisterFileTypeDecoder_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterFileTypeDecoder[BadFile](badFileDecoder{}, true)
	})

	assert.NotPanics(t, func() {
		RegisterFileTypeDecoder[BadFile](badFileDecoder{}, true)
	})

	removeTypeDecoder[BadFile]() // remove the custom decoder
}

func TestRegisterNamedDecoder(t *testing.T) {
	assert.Panics(t, func() { RegisterNamedDecoder[bool]("myBool", nil) }) // fail on nil decoder

	assert.Panics(t, func() {
		// Fail on invalid decoder (invalid signature).
		RegisterNamedDecoder[bool]("myBool", func(string) error {
			return nil
		})
	})

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

func Test_smartDecoder_BasicTypes(t *testing.T) {
	// returns int
	intDecoder := DecoderFunc[string](decodeInt)

	// returns *int
	intPointerDecoder := DecoderFunc[string](func(value string) (interface{}, error) {
		if v, err := decodeInt(value); err != nil {
			return nil, err
		} else {
			var x = v.(int)
			return &x, nil
		}
	})

	intType := typeOf[int]()
	intPointerType := typeOf[*int]()

	smartIntDecoders := []ValueTypeDecoder{
		newSmartDecoderX(intType, intDecoder).(ValueTypeDecoder),
		newSmartDecoderX(intType, intPointerDecoder).(ValueTypeDecoder),
	}
	smartIntPointerDecoders := []ValueTypeDecoder{
		newSmartDecoderX(intPointerType, intDecoder).(ValueTypeDecoder),
		newSmartDecoderX(intPointerType, intPointerDecoder).(ValueTypeDecoder),
	}

	for _, decoder := range smartIntDecoders {
		v, err := decoder.Decode("2000")
		success[int](t, 2000, v, err)
	}

	for _, decoder := range smartIntPointerDecoders {
		v, err := decoder.Decode("2000")
		var ev int = 2000
		success[*int](t, &ev, v, err)
	}
}

func Test_smartDecoder_StructTypes(t *testing.T) {
	placeType := typeOf[Place]()
	placePointerType := typeOf[*Place]()

	// myPlaceDecoder returns Place
	// myPlacePointerDecoder returns *Place

	smartPlaceDecoders := []ValueTypeDecoder{
		newSmartDecoderX(placeType, myPlaceDecoder).(ValueTypeDecoder),
		newSmartDecoderX(placeType, myPlacePointerDecoder).(ValueTypeDecoder),
	}

	smartPlacePointerDecoders := []ValueTypeDecoder{
		newSmartDecoderX(placePointerType, myPlaceDecoder).(ValueTypeDecoder),
		newSmartDecoderX(placePointerType, myPlacePointerDecoder).(ValueTypeDecoder),
	}

	for _, decoder := range smartPlaceDecoders {
		v, err := decoder.Decode("Canada.Toronto")
		success[Place](t, Place{Country: "Canada", City: "Toronto"}, v, err)
	}

	for _, decoder := range smartPlacePointerDecoders {
		v, err := decoder.Decode("Canada.Toronto")
		success[*Place](t, &Place{Country: "Canada", City: "Toronto"}, v, err)
	}
}

func Test_smartDecoder_ErrValueTypeMismatch(t *testing.T) {
	// myDateDecoder decodes a string to a time.Time.
	// While we set the desired type to int, so it should fail.
	smart := newSmartDecoder[string](typeOf[int](), myDateDecoder)
	v, err := smart.Decode("2001-02-03")
	assert.Nil(t, v)
	assert.ErrorIs(t, err, ErrTypeMismatch)
	assert.ErrorContains(t, err, invalidDecodeReturnType(reflect.TypeOf(0), reflect.TypeOf(time.Time{})).Error())
}

// Test that the builtin decoders are valid.

func TestDecoder_bool(t *testing.T) {
	v, err := decodeBool("true")
	success[bool](t, true, v, err)

	v, err = decodeBool("false")
	success[bool](t, false, v, err)

	v, err = decodeBool("1")
	success[bool](t, true, v, err)

	v, err = decodeBool("apple")
	fail[bool](t, false, v, err)
}

func TestDecoder_int(t *testing.T) {
	v, err := decodeInt("2045")
	success[int](t, 2045, v, err)

	v, err = decodeInt("apple")
	fail[int](t, 0, v, err)
}

func TestDecoder_int8(t *testing.T) {
	v, err := decodeInt8("127")
	success[int8](t, 127, v, err)

	v, err = decodeInt8("128")
	fail[int8](t, 127, v, err)

	v, err = decodeInt8("apple")
	fail[int8](t, 0, v, err)
}

func TestDecoder_int16(t *testing.T) {
	v, err := decodeInt16("32767")
	success[int16](t, 32767, v, err)

	v, err = decodeInt16("32768")
	fail[int16](t, 32767, v, err)

	v, err = decodeInt16("apple")
	fail[int16](t, 0, v, err)
}

func TestDecoder_int32(t *testing.T) {
	v, err := decodeInt32("2147483647")
	success[int32](t, 2147483647, v, err)

	v, err = decodeInt32("2147483648")
	fail[int32](t, 2147483647, v, err)

	v, err = decodeInt32("apple")
	fail[int32](t, 0, v, err)
}

func TestDecoder_int64(t *testing.T) {
	v, err := decodeInt64("9223372036854775807")
	success[int64](t, 9223372036854775807, v, err)

	v, err = decodeInt64("9223372036854775808")
	fail[int64](t, 9223372036854775807, v, err)

	v, err = decodeInt64("apple")
	fail[int64](t, 0, v, err)
}

func TestDecoder_uint(t *testing.T) {
	v, err := decodeUint("2045")
	success[uint](t, 2045, v, err)

	v, err = decodeUint("apple")
	fail[uint](t, 0, v, err)
}

func TestDecoder_uint8(t *testing.T) {
	v, err := decodeUint8("255")
	success[uint8](t, 255, v, err)

	v, err = decodeUint8("256")
	fail[uint8](t, 255, v, err)

	v, err = decodeUint8("apple")
	fail[uint8](t, 0, v, err)
}

func TestDecoder_uint16(t *testing.T) {
	v, err := decodeUint16("65535")
	success[uint16](t, 65535, v, err)

	v, err = decodeUint16("65536")
	fail[uint16](t, 65535, v, err)

	v, err = decodeUint16("apple")
	fail[uint16](t, 0, v, err)
}

func TestDecoder_uint32(t *testing.T) {
	v, err := decodeUint32("4294967295")
	success[uint32](t, 4294967295, v, err)

	v, err = decodeUint32("4294967296")
	fail[uint32](t, 4294967295, v, err)

	v, err = decodeUint32("apple")
	fail[uint32](t, 0, v, err)
}

func TestDecoder_uint64(t *testing.T) {
	v, err := decodeUint64("18446744073709551615")
	success[uint64](t, 18446744073709551615, v, err)

	v, err = decodeUint64("18446744073709551616")
	fail[uint64](t, 18446744073709551615, v, err)

	v, err = decodeUint64("apple")
	fail[uint64](t, 0, v, err)
}

func TestDecoder_float32(t *testing.T) {
	v, err := decodeFloat32("3.1415926")
	success[float32](t, 3.1415926, v, err)

	v, err = decodeFloat32("apple")
	fail[float32](t, 0, v, err)
}

func TestDecoder_float64(t *testing.T) {
	v, err := decodeFloat64("3.1415926")
	success[float64](t, 3.1415926, v, err)

	v, err = decodeFloat64("apple")
	fail[float64](t, 0, v, err)
}

func TestDecoder_complex64(t *testing.T) {
	v, err := decodeComplex64("1+4i")
	success[complex64](t, 1+4i, v, err)

	v, err = decodeComplex64("apple")
	fail[complex64](t, 0, v, err)
}

func TestDecoder_complex128(t *testing.T) {
	v, err := decodeComplex128("1+4i")
	success[complex128](t, 1+4i, v, err)

	v, err = decodeComplex128("apple")
	fail[complex128](t, 0, v, err)
}

func TestDecoder_string(t *testing.T) {
	v, err := decodeString("hello")
	success[string](t, "hello", v, err)
}

func TestDecoder_time(t *testing.T) {
	v, err := decodeTime("1991-11-10T08:00:00+08:00")
	assert.NoError(t, err)
	expected := time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("Asia/Shanghai", +8*3600))
	assert.True(t, equalTime(expected, v.(time.Time)))

	v, err = decodeTime("678088800")
	expected = time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC)
	assert.NoError(t, err)
	assert.True(t, equalTime(expected, v.(time.Time)))

	v, err = decodeTime("678088800.123456")
	assert.NoError(t, err)
	expected = time.Date(1991, 6, 28, 6, 0, 0, 123456000, time.UTC)
	assert.True(t, equalTime(expected, v.(time.Time)))

	v, err = decodeTime("apple")
	assert.Error(t, err)
	assert.True(t, equalTime(time.Time{}, v.(time.Time)))
}

func equalFuncs(expected, actual interface{}) bool {
	return reflect.ValueOf(expected).Pointer() == reflect.ValueOf(actual).Pointer()
}

func equalTime(expected, actual time.Time) bool {
	return expected.UTC() == actual.UTC()
}

func removeTypeDecoder[T any]() {
	delete(customDecoders, typeOf[T]())
	delete(customDecoders, typeOf[[]T]())
	delete(customDecoders, typeOf[patch.Field[T]]())
	delete(customDecoders, typeOf[patch.Field[[]T]]())
}

func removeNamedDecoder(name string) {
	delete(namedDecoders, name)
}

func success[T any](t *testing.T, expected T, got interface{}, err error) {
	assert.NoError(t, err)
	_, ok := got.(T)
	assert.True(t, ok)
	assert.Equal(t, expected, got)
}

func fail[T any](t *testing.T, expected T, got interface{}, err error) {
	assert.Error(t, err)
	_, ok := got.(T)
	assert.True(t, ok)
	assert.Equal(t, expected, got)
}
