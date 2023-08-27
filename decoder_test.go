package httpin

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

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

func invalidDecoder(string) error {
	return nil
}

func TestRegisterTypeDecoder(t *testing.T) {
	boolType := reflect.TypeOf(bool(true))
	assert.Panics(t, func() { RegisterTypeDecoder(boolType, nil) })
	assert.Panics(t, func() { RegisterTypeDecoder(boolType, invalidDecoder) })
	delete(decoders, boolType) // remove the custom decoder

	// Register duplicate decoder should fail.
	assert.NotPanics(t, func() {
		RegisterTypeDecoder(boolType, ValueTypeDecoderFunc(decodeCustomBool))
	})
	assert.Panics(t, func() {
		RegisterTypeDecoder(boolType, ValueTypeDecoderFunc(decodeCustomBool))
	})

	delete(decoders, boolType) // remove the custom decoder
}

func TestRegisterNamedDecoder(t *testing.T) {
	assert.Panics(t, func() { RegisterNamedDecoder("myBool", nil) })
	assert.Panics(t, func() { RegisterNamedDecoder("myBool", invalidDecoder) })
	delete(namedDecoders, "mybool") // remove the custom decoder

	// Register duplicate decoder should fail.
	assert.NotPanics(t, func() {
		RegisterNamedDecoder("mybool", ValueTypeDecoderFunc(decodeCustomBool))
	})
	assert.Panics(t, func() {
		RegisterNamedDecoder("mybool", ValueTypeDecoderFunc(decodeCustomBool))
	})

	delete(namedDecoders, "mybool") // remove the custom decoder
}

func TestReplaceTypeDecoder(t *testing.T) {
	boolType := reflect.TypeOf(bool(true))

	assert.NotPanics(t, func() {
		ReplaceTypeDecoder(boolType, ValueTypeDecoderFunc(decodeCustomBool))
	})

	assert.NotPanics(t, func() {
		ReplaceTypeDecoder(boolType, ValueTypeDecoderFunc(decodeCustomBool))
	})
}

func TestReplaceNamedDecoder(t *testing.T) {
	assert.NotPanics(t, func() {
		ReplaceNamedDecoder("mybool", ValueTypeDecoderFunc(decodeCustomBool))
	})

	assert.NotPanics(t, func() {
		ReplaceNamedDecoder("mybool", ValueTypeDecoderFunc(decodeCustomBool))
	})
}
