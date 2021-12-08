package httpin

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// DecodeCustomBool additionally parses "yes/no".
func DecodeCustomBool(value string) (interface{}, error) {
	sdata := strings.ToLower(value)
	if sdata == "yes" {
		return true, nil
	}
	if sdata == "no" {
		return false, nil
	}
	return strconv.ParseBool(sdata)
}

func TestDecoders(t *testing.T) {
	boolType := reflect.TypeOf(bool(true))

	Convey("Register nil decoder", t, func() {
		So(func() { RegisterTypeDecoder(boolType, nil) }, ShouldPanic)
	})
	delete(decoders, boolType) // remove the custom decoder

	var invalidDecoder = func(string) error {
		return nil
	}

	Convey("Register invalid decoder", t, func() {
		So(func() { RegisterTypeDecoder(boolType, invalidDecoder) }, ShouldPanic)
	})
	delete(decoders, boolType) // remove the custom decoder

	Convey("Register duplicate decoder", t, func() {
		So(func() { RegisterTypeDecoder(boolType, ValueTypeDecoderFunc(DecodeCustomBool)) }, ShouldNotPanic)
		So(func() { RegisterTypeDecoder(boolType, ValueTypeDecoderFunc(DecodeCustomBool)) }, ShouldPanic)
	})
	delete(decoders, boolType) // remove the custom decoder

	Convey("Replace a decoder", t, func() {
		So(func() { ReplaceTypeDecoder(boolType, ValueTypeDecoderFunc(DecodeCustomBool)) }, ShouldNotPanic)
		So(func() { ReplaceTypeDecoder(boolType, ValueTypeDecoderFunc(DecodeCustomBool)) }, ShouldNotPanic)
	})
	delete(decoders, boolType) // remove the custom decoder
}
