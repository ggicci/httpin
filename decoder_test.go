package httpin

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// DecodeCustomBool additionally parses "yes/no".
func DecodeCustomBool(data []byte) (interface{}, error) {
	sdata := strings.ToLower(string(data))
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
		So(func() { RegisterDecoder(boolType, nil) }, ShouldPanic)
	})
	delete(decoders, boolType) // remove the custom decoder

	Convey("Register duplicate decoder", t, func() {
		So(func() { RegisterDecoder(boolType, DecoderFunc(DecodeCustomBool)) }, ShouldNotPanic)
		So(func() { RegisterDecoder(boolType, DecoderFunc(DecodeCustomBool)) }, ShouldPanic)
	})
	delete(decoders, boolType) // remove the custom decoder

	Convey("Replace a decoder", t, func() {
		So(func() { ReplaceDecoder(boolType, DecoderFunc(DecodeCustomBool)) }, ShouldNotPanic)
		So(func() { ReplaceDecoder(boolType, DecoderFunc(DecodeCustomBool)) }, ShouldNotPanic)
	})
	delete(decoders, boolType) // remove the custom decoder
}
