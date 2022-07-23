package internal

import (
	"mime/multipart"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type Thing struct{}

func TestBuiltinDecoders(t *testing.T) {

	Convey("DecoderFunc implements Decoder interface", t, func() {
		v, err := ValueTypeDecoderFunc(DecodeBool).Decode("true")
		So(v, ShouldBeTrue)
		So(err, ShouldBeNil)
	})

	Convey("DecoderOf retrieves a decoder by type", t, func() {
		So(DecoderOf(reflect.TypeOf(true)), ShouldEqual, DecodeBool)
		So(DecoderOf(reflect.TypeOf(Thing{})), ShouldBeNil)
	})

	Convey("Decoder for bool type", t, func() {
		v, err := DecodeBool("true")
		So(v, ShouldBeTrue)
		So(v, ShouldHaveSameTypeAs, true)
		So(err, ShouldBeNil)
		v, err = DecodeBool("false")
		So(v, ShouldBeFalse)
		So(err, ShouldBeNil)
		v, err = DecodeBool("1")
		So(v, ShouldBeTrue)
		So(err, ShouldBeNil)
		_, err = DecodeBool("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for int type", t, func() {
		v, err := DecodeInt("2045")
		So(v, ShouldEqual, 2045)
		So(v, ShouldHaveSameTypeAs, int(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for int8 type", t, func() {
		v, err := DecodeInt8("127")
		So(v, ShouldEqual, 127)
		So(v, ShouldHaveSameTypeAs, int8(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt8("128")
		So(err, ShouldBeError)
		_, err = DecodeInt8("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for int16 type", t, func() {
		v, err := DecodeInt16("32767")
		So(v, ShouldEqual, 32767)
		So(v, ShouldHaveSameTypeAs, int16(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt16("32768")
		So(err, ShouldBeError)
		_, err = DecodeInt16("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for int32 type", t, func() {
		v, err := DecodeInt32("2147483647")
		So(v, ShouldEqual, 2147483647)
		So(v, ShouldHaveSameTypeAs, int32(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt32("2147483648")
		So(err, ShouldBeError)
		_, err = DecodeInt32("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for int64 type", t, func() {
		v, err := DecodeInt64("9223372036854775807")
		So(v, ShouldEqual, 9223372036854775807)
		So(v, ShouldHaveSameTypeAs, int64(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt64("9223372036854775808")
		So(err, ShouldBeError)
		_, err = DecodeInt64("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint type", t, func() {
		v, err := DecodeUint("2045")
		So(v, ShouldEqual, uint(2045))
		So(v, ShouldHaveSameTypeAs, uint(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint8 type", t, func() {
		v, err := DecodeUint8("255")
		So(v, ShouldEqual, uint8(255))
		So(v, ShouldHaveSameTypeAs, uint8(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint8("256")
		So(err, ShouldBeError)
		_, err = DecodeUint8("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint16 type", t, func() {
		v, err := DecodeUint16("65535")
		So(v, ShouldEqual, uint16(65535))
		So(v, ShouldHaveSameTypeAs, uint16(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint16("65536")
		So(err, ShouldBeError)
		_, err = DecodeUint16("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint32 type", t, func() {
		v, err := DecodeUint32("4294967295")
		So(v, ShouldEqual, uint32(4294967295))
		So(v, ShouldHaveSameTypeAs, uint32(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint32("4294967296")
		So(err, ShouldBeError)
		_, err = DecodeUint32("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint64 type", t, func() {
		v, err := DecodeUint64("18446744073709551615")
		So(v, ShouldEqual, uint64(18446744073709551615))
		So(v, ShouldHaveSameTypeAs, uint64(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint64("18446744073709551616")
		So(err, ShouldBeError)
		_, err = DecodeUint64("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for float32 type", t, func() {
		v, err := DecodeFloat32("3.1415926")
		So(v, ShouldEqual, 3.1415926)
		So(v, ShouldHaveSameTypeAs, float32(0.0))
		So(err, ShouldBeNil)
		_, err = DecodeFloat32("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for float64 type", t, func() {
		v, err := DecodeFloat64("3.1415926")
		So(v, ShouldEqual, 3.1415926)
		So(v, ShouldHaveSameTypeAs, float64(0.0))
		So(err, ShouldBeNil)
		_, err = DecodeFloat64("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for complex64 type", t, func() {
		v, err := DecodeComplex64("1+4i")
		So(v, ShouldEqual, 1+4i)
		So(v, ShouldHaveSameTypeAs, complex64(1+4i))
		So(err, ShouldBeNil)
		_, err = DecodeComplex64("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for complex128 type", t, func() {
		v, err := DecodeComplex128("1+4i")
		So(v, ShouldEqual, 1+4i)
		So(v, ShouldHaveSameTypeAs, complex128(1+4i))
		So(err, ShouldBeNil)
		_, err = DecodeComplex128("apple")
		So(err, ShouldBeError)
	})

	Convey("Decoder for string type", t, func() {
		v, err := DecodeString("hello")
		So(v, ShouldEqual, "hello")
		So(v, ShouldHaveSameTypeAs, string(""))
		So(err, ShouldBeNil)
	})

	Convey("Decoder for time.Time type", t, func() {
		v, err := DecodeTime("1991-11-10T08:00:00+08:00")
		So(v, ShouldEqual, time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("Asia/Shanghai", +8*3600)))
		So(v, ShouldHaveSameTypeAs, time.Time{})
		So(v.(time.Time).Location(), ShouldEqual, time.UTC)
		So(err, ShouldBeNil)

		v, err = DecodeTime("678088800")
		So(v, ShouldEqual, time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC))
		So(v, ShouldHaveSameTypeAs, time.Time{})
		So(v.(time.Time).Location(), ShouldEqual, time.UTC)
		So(err, ShouldBeNil)

		v, err = DecodeTime("678088800.123456")
		So(v, ShouldEqual, time.Date(1991, 6, 28, 6, 0, 0, 123456000, time.UTC))
		So(v, ShouldHaveSameTypeAs, time.Time{})
		So(v.(time.Time).Location(), ShouldEqual, time.UTC)
		So(err, ShouldBeNil)

		_, err = DecodeTime("apple")
		So(err, ShouldBeError)
	})
}

func TestTypeDecoderAdapter(t *testing.T) {
	Convey("Adapter: ValueTypeDecoderFunc", t, func() {

		decoder := ValueTypeDecoderFunc(func(value string) (interface{}, error) {
			return value + "!", nil
		})

		got, err := decoder.Decode("hello")
		So(err, ShouldBeNil)
		So(got, ShouldEqual, "hello!")
	})

	Convey("Adapter: FileTypeDecoderFunc", t, func() {
		decoder := FileTypeDecoderFunc(func(file *multipart.FileHeader) (interface{}, error) {
			return file.Filename, nil
		})

		fileHeader := &multipart.FileHeader{
			Filename: "hello.txt",
		}

		got, err := decoder.Decode(fileHeader)
		So(err, ShouldBeNil)
		So(got, ShouldEqual, "hello.txt")
	})
}
