package httpin_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ggicci/httpin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTypes_Bool(t *testing.T) {
	var v httpin.Bool

	Convey("Marshal Bool", t, func() {
		So(json.Unmarshal([]byte("true"), &v), ShouldBeNil)
		So(v.Value, ShouldBeTrue)
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Unmarshal Bool", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("true"))
	})
}

func TestTypes_Int(t *testing.T) {
	var v httpin.Int

	Convey("Marshal Int", t, func() {
		So(json.Unmarshal([]byte("2015"), &v), ShouldBeNil)
		So(v.Value, ShouldEqual, 2015)
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Unmarshal Int", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("2015"))
	})
}

func TestTypes_Uint(t *testing.T) {
	var v httpin.Uint

	Convey("Marshal Int", t, func() {
		So(json.Unmarshal([]byte("2045"), &v), ShouldBeNil)
		So(v.Value, ShouldEqual, 2045)
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Unmarshal Int", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("2045"))
	})
}

func TestTypes_Float32(t *testing.T) {
	var v httpin.Float32

	Convey("Marshal Float32", t, func() {
		So(json.Unmarshal([]byte("3.1415"), &v), ShouldBeNil)
		So(v.Value, ShouldEqual, 3.1415)
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Unmarshal Float32", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("3.1415"))
	})
}

func TestTypes_Time(t *testing.T) {
	var v httpin.Time

	Convey("Marshal Time", t, func() {
		So(json.Unmarshal([]byte("1991-11-10T08:00:00+08:00"), &v), ShouldBeNil)
		So(v.Value, ShouldEqual, time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("E8", 8*3600)))
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Unmarshal Time", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("1991-11-10T08:00:00+08:00"))
	})
}

func TestTypes_BoolArray(t *testing.T) {
	var v httpin.BoolArray

	Convey("Marshal BoolArray", t, func() {
		So(json.Unmarshal([]byte("[true, false, true]"), &v), ShouldBeNil)
		So(v.Value, ShouldResemble, []bool{true, false, true})
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Unmarshal BoolArray", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("[true,false,true]"))
	})
}

func TestTypes_IntArray(t *testing.T) {
	var v httpin.IntArray

	Convey("Marshal IntArray", t, func() {
		So(json.Unmarshal([]byte("[9, 12, 1024]"), &v), ShouldBeNil)
		So(v.Value, ShouldResemble, []int{9, 12, 1024})
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Unmarshal IntArray", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("[9,12,1024]"))
	})
}
