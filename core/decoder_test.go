package core

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ggicci/httpin/internal"
	"github.com/stretchr/testify/assert"
)

func decodeMyDate(value string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, &InvalidDate{Value: value, Err: err}
	}
	return t, nil
}

var myDateDecoder = DecoderFunc[time.Time](decodeMyDate)

type Place struct {
	Country string
	City    string
}

func (p Place) ToString() (string, error) {
	return fmt.Sprintf("%s.%s", p.Country, p.City), nil
}

func (p *Place) FromString(value string) error {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return errors.New("invalid place")
	}
	*p = Place{Country: parts[0], City: parts[1]}
	return nil
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

func decodeInt(value string) (int, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	return int(v), err
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

	intType := internal.TypeOf[int]()
	intPointerType := internal.TypeOf[*int]()

	smartIntDecoders := []Decoder[any]{
		NewSmartDecoder(intType, ToAnyDecoder[int](intDecoder)),
		NewSmartDecoder(intType, ToAnyDecoder[*int](intPointerDecoder)),
	}
	smartIntPointerDecoders := []Decoder[any]{
		NewSmartDecoder(intPointerType, ToAnyDecoder[int](intDecoder)),
		NewSmartDecoder(intPointerType, ToAnyDecoder[*int](intPointerDecoder)),
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
	placeType := internal.TypeOf[Place]()
	placePointerType := internal.TypeOf[*Place]()

	// myPlaceDecoder returns Place
	// myPlacePointerDecoder returns *Place

	smartPlaceDecoders := []Decoder[any]{
		NewSmartDecoder(placeType, ToAnyDecoder[Place](myPlaceDecoder)),
		NewSmartDecoder(placeType, ToAnyDecoder[*Place](myPlacePointerDecoder)),
	}

	smartPlacePointerDecoders := []Decoder[any]{
		NewSmartDecoder(placePointerType, ToAnyDecoder[Place](myPlaceDecoder)),
		NewSmartDecoder(placePointerType, ToAnyDecoder[*Place](myPlacePointerDecoder)),
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
	decoder := ToAnyDecoder[*int](DecoderFunc[*int](func(value string) (*int, error) {
		return nil, nil
	}))
	v, err := NewSmartDecoder(internal.TypeOf[int](), decoder).Decode("100")
	assert.Nil(t, v)
	assert.NoError(t, err)
}

func TestSmartDecoder_ErrValueTypeMismatch(t *testing.T) {
	// myDateDecoder decodes a string to a time.Time.
	// While we set the desired type to int, so it should fail.
	smart := NewSmartDecoder(internal.TypeOf[int](), ToAnyDecoder[time.Time](myDateDecoder))
	v, err := smart.Decode("2001-02-03")
	assert.Nil(t, v)
	assert.ErrorIs(t, err, errTypeMismatch)
	assert.ErrorContains(t, err, typeMismatchedError(reflect.TypeOf(0), reflect.TypeOf(time.Time{})).Error())
}

func success[T any](t *testing.T, expected T, got any, err error) {
	assert.NoError(t, err)
	_, ok := got.(T)
	assert.True(t, ok)
	assert.Equal(t, expected, got)
}
