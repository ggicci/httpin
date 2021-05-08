package internal

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type Thing struct{}

func TestBuiltinDecoders(t *testing.T) {

	Convey("DecoderFunc implements Decoder interface", t, func() {
		v, err := DecoderFunc(DecodeBool).Decode([]byte("true"))
		So(v, ShouldBeTrue)
		So(err, ShouldBeNil)
	})

	Convey("DecoderOf retrieves a decoder by type", t, func() {
		So(DecoderOf(reflect.TypeOf(true)), ShouldEqual, DecodeBool)
		So(DecoderOf(reflect.TypeOf(Thing{})), ShouldBeNil)
	})

	Convey("Decoder for bool type", t, func() {
		v, err := DecodeBool([]byte("true"))
		So(v, ShouldBeTrue)
		So(v, ShouldHaveSameTypeAs, true)
		So(err, ShouldBeNil)
		v, err = DecodeBool([]byte("false"))
		So(v, ShouldBeFalse)
		So(err, ShouldBeNil)
		v, err = DecodeBool([]byte("1"))
		So(v, ShouldBeTrue)
		So(err, ShouldBeNil)
		_, err = DecodeBool([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for int type", t, func() {
		v, err := DecodeInt([]byte("2045"))
		So(v, ShouldEqual, 2045)
		So(v, ShouldHaveSameTypeAs, int(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for int8 type", t, func() {
		v, err := DecodeInt8([]byte("127"))
		So(v, ShouldEqual, 127)
		So(v, ShouldHaveSameTypeAs, int8(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt8([]byte("128"))
		So(err, ShouldBeError)
		_, err = DecodeInt8([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for int16 type", t, func() {
		v, err := DecodeInt16([]byte("32767"))
		So(v, ShouldEqual, 32767)
		So(v, ShouldHaveSameTypeAs, int16(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt16([]byte("32768"))
		So(err, ShouldBeError)
		_, err = DecodeInt16([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for int32 type", t, func() {
		v, err := DecodeInt32([]byte("2147483647"))
		So(v, ShouldEqual, 2147483647)
		So(v, ShouldHaveSameTypeAs, int32(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt32([]byte("2147483648"))
		So(err, ShouldBeError)
		_, err = DecodeInt32([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for int64 type", t, func() {
		v, err := DecodeInt64([]byte("9223372036854775807"))
		So(v, ShouldEqual, 9223372036854775807)
		So(v, ShouldHaveSameTypeAs, int64(1))
		So(err, ShouldBeNil)
		_, err = DecodeInt64([]byte("9223372036854775808"))
		So(err, ShouldBeError)
		_, err = DecodeInt64([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint type", t, func() {
		v, err := DecodeUint([]byte("2045"))
		So(v, ShouldEqual, uint(2045))
		So(v, ShouldHaveSameTypeAs, uint(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint8 type", t, func() {
		v, err := DecodeUint8([]byte("255"))
		So(v, ShouldEqual, uint8(255))
		So(v, ShouldHaveSameTypeAs, uint8(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint8([]byte("256"))
		So(err, ShouldBeError)
		_, err = DecodeUint8([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint16 type", t, func() {
		v, err := DecodeUint16([]byte("65535"))
		So(v, ShouldEqual, uint16(65535))
		So(v, ShouldHaveSameTypeAs, uint16(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint16([]byte("65536"))
		So(err, ShouldBeError)
		_, err = DecodeUint16([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint32 type", t, func() {
		v, err := DecodeUint32([]byte("4294967295"))
		So(v, ShouldEqual, uint32(4294967295))
		So(v, ShouldHaveSameTypeAs, uint32(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint32([]byte("4294967296"))
		So(err, ShouldBeError)
		_, err = DecodeUint32([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for uint64 type", t, func() {
		v, err := DecodeUint64([]byte("18446744073709551615"))
		So(v, ShouldEqual, uint64(18446744073709551615))
		So(v, ShouldHaveSameTypeAs, uint64(1))
		So(err, ShouldBeNil)
		_, err = DecodeUint64([]byte("18446744073709551616"))
		So(err, ShouldBeError)
		_, err = DecodeUint64([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for float32 type", t, func() {
		v, err := DecodeFloat32([]byte("3.1415926"))
		So(v, ShouldEqual, 3.1415926)
		So(v, ShouldHaveSameTypeAs, float32(0.0))
		So(err, ShouldBeNil)
		_, err = DecodeFloat32([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for float64 type", t, func() {
		v, err := DecodeFloat64([]byte("3.1415926"))
		So(v, ShouldEqual, 3.1415926)
		So(v, ShouldHaveSameTypeAs, float64(0.0))
		So(err, ShouldBeNil)
		_, err = DecodeFloat64([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for complex64 type", t, func() {
		v, err := DecodeComplex64([]byte("1+4i"))
		So(v, ShouldEqual, 1+4i)
		So(v, ShouldHaveSameTypeAs, complex64(1+4i))
		So(err, ShouldBeNil)
		_, err = DecodeComplex64([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for complex128 type", t, func() {
		v, err := DecodeComplex128([]byte("1+4i"))
		So(v, ShouldEqual, 1+4i)
		So(v, ShouldHaveSameTypeAs, complex128(1+4i))
		So(err, ShouldBeNil)
		_, err = DecodeComplex128([]byte("apple"))
		So(err, ShouldBeError)
	})

	Convey("Decoder for string type", t, func() {
		v, err := DecodeString([]byte("hello"))
		So(v, ShouldEqual, "hello")
		So(v, ShouldHaveSameTypeAs, string(""))
		So(err, ShouldBeNil)
	})

	Convey("Decoder for time.Time type", t, func() {
		v, err := DecodeTime([]byte("1991-11-10T08:00:00+08:00"))
		So(v, ShouldEqual, time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("Asia/Shanghai", +8*3600)))
		So(v, ShouldHaveSameTypeAs, time.Time{})
		So(v.(time.Time).Location(), ShouldEqual, time.UTC)
		So(err, ShouldBeNil)

		v, err = DecodeTime([]byte("678088800"))
		So(v, ShouldEqual, time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC))
		So(v, ShouldHaveSameTypeAs, time.Time{})
		So(v.(time.Time).Location(), ShouldEqual, time.UTC)
		So(err, ShouldBeNil)

		_, err = DecodeTime([]byte("apple"))
		So(err, ShouldBeError)
	})
}
