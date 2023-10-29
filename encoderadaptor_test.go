package httpin

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoderAdaptor(t *testing.T) {

}

type Location struct {
	Latitude  float64
	Longitude float64
}

func (l Location) String() string {
	return fmt.Sprintf("%f,%f", l.Latitude, l.Longitude)
}

type LocationFormValueMarshalerImpl Location

func (l LocationFormValueMarshalerImpl) HttpinFormValue() (string, error) {
	return "HttpinFormValue:" + (Location)(l).String(), nil
}

type LocationTextMarshalerImpl Location

func (l LocationTextMarshalerImpl) MarshalText() ([]byte, error) {
	return []byte("MarshalText:" + (Location)(l).String()), nil
}

func TestInterfaceEncoder_fallbackEncoder(t *testing.T) {
	loc := &Location{
		Latitude:  1.234,
		Longitude: 5.678,
	}

	actual, err := fallbackEncoder.Encode(reflect.ValueOf(loc))
	assert.NoError(t, err)
	assert.Equal(t, "1.234000,5.678000", actual)

	actual, err = fallbackEncoder.Encode(reflect.ValueOf(LocationFormValueMarshalerImpl(*loc)))
	assert.NoError(t, err)
	assert.Equal(t, "HttpinFormValue:1.234000,5.678000", actual)

	actual, err = fallbackEncoder.Encode(reflect.ValueOf(LocationTextMarshalerImpl(*loc)))
	assert.NoError(t, err)
	assert.Equal(t, "MarshalText:1.234000,5.678000", actual)
}
