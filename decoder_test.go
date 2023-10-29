package httpin

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

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

type Place struct {
	Country string
	City    string
}

// decodePlace parses "country.city", e.g. "Canada.Toronto".
// It returns a Place.
func decodePlace(value string) (Place, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return Place{}, errors.New("invalid place")
	}
	return Place{Country: parts[0], City: parts[1]}, nil
}

// decodePlacePointer parses "country.city", e.g. "Canada.Toronto".
// It returns *Place.
func decodePlacePointer(value string) (*Place, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid place")
	}
	return &Place{Country: parts[0], City: parts[1]}, nil
}

var myPlaceDecoder = DecoderFunc[Place](decodePlace)
var myPlacePointerDecoder = DecoderFunc[*Place](decodePlacePointer)

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

func TestSmartDecoder_BasicTypes(t *testing.T) {
	// returns int
	intDecoder := DecoderFunc[int](decodeInt)

	// returns *int
	intPointerDecoder := DecoderFunc[*int](func(value string) (*int, error) {
		if v, err := decodeInt(value); err != nil {
			return nil, err
		} else {
			var x = v
			return &x, nil
		}
	})

	intType := typeOf[int]()
	intPointerType := typeOf[*int]()

	smartIntDecoders := []Decoder[any]{
		newSmartDecoder(intType, toAnyDecoder[int](intDecoder)),
		newSmartDecoder(intType, toAnyDecoder[*int](intPointerDecoder)),
	}
	smartIntPointerDecoders := []Decoder[any]{
		newSmartDecoder(intPointerType, toAnyDecoder[int](intDecoder)),
		newSmartDecoder(intPointerType, toAnyDecoder[*int](intPointerDecoder)),
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

func TestSmartDecoder_StructTypes(t *testing.T) {
	placeType := typeOf[Place]()
	placePointerType := typeOf[*Place]()

	// myPlaceDecoder returns Place
	// myPlacePointerDecoder returns *Place

	smartPlaceDecoders := []Decoder[any]{
		newSmartDecoder(placeType, toAnyDecoder[Place](myPlaceDecoder)),
		newSmartDecoder(placeType, toAnyDecoder[*Place](myPlacePointerDecoder)),
	}

	smartPlacePointerDecoders := []Decoder[any]{
		newSmartDecoder(placePointerType, toAnyDecoder[Place](myPlaceDecoder)),
		newSmartDecoder(placePointerType, toAnyDecoder[*Place](myPlacePointerDecoder)),
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

func TestSmartDecoder_ReturnNil(t *testing.T) {
	decoder := toAnyDecoder[*int](DecoderFunc[*int](func(value string) (*int, error) {
		return nil, nil
	}))
	v, err := newSmartDecoder(typeOf[int](), decoder).Decode("100")
	assert.Nil(t, v)
	assert.NoError(t, err)
}

func TestSmartDecoder_ErrValueTypeMismatch(t *testing.T) {
	// myDateDecoder decodes a string to a time.Time.
	// While we set the desired type to int, so it should fail.
	smart := newSmartDecoder(typeOf[int](), toAnyDecoder[time.Time](myDateDecoder))
	v, err := smart.Decode("2001-02-03")
	assert.Nil(t, v)
	assert.ErrorIs(t, err, errTypeMismatch)
	assert.ErrorContains(t, err, invalidDecodeReturnType(reflect.TypeOf(0), reflect.TypeOf(time.Time{})).Error())
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
