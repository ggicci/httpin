package httpin_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ggicci/httpin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTypes_Bool(t *testing.T) {
	var v httpin.Bool

	Convey("Unmarshal Bool", t, func() {
		So(json.Unmarshal([]byte("true"), &v), ShouldBeNil)
		So(v.Value, ShouldBeTrue)
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal Bool", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("true"))
	})
}

func TestTypes_Int(t *testing.T) {
	var v httpin.Int

	Convey("Unmarshal Int", t, func() {
		So(json.Unmarshal([]byte("2015"), &v), ShouldBeNil)
		So(v.Value, ShouldEqual, 2015)
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal Int", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("2015"))
	})
}

func TestTypes_Uint(t *testing.T) {
	var v httpin.Uint

	Convey("Unmarshal Int", t, func() {
		So(json.Unmarshal([]byte("2045"), &v), ShouldBeNil)
		So(v.Value, ShouldEqual, 2045)
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal Int", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("2045"))
	})
}

func TestTypes_Float32(t *testing.T) {
	var v httpin.Float32

	Convey("Unmarshal Float32", t, func() {
		So(json.Unmarshal([]byte("3.1415"), &v), ShouldBeNil)
		So(v.Value, ShouldEqual, 3.1415)
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal Float32", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("3.1415"))
	})
}

func TestTypes_Time(t *testing.T) {
	var v httpin.Time

	Convey("Unmarshal Time", t, func() {
		So(json.Unmarshal([]byte("\"1991-11-10T08:00:00+08:00\""), &v), ShouldBeNil)
		So(v.Value, ShouldEqual, time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("E8", 8*3600)))
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal Time", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("\"1991-11-10T08:00:00+08:00\""))
	})
}

func TestTypes_BoolArray(t *testing.T) {
	var v httpin.BoolArray

	Convey("Unmarshal BoolArray", t, func() {
		So(json.Unmarshal([]byte("[true, false, true]"), &v), ShouldBeNil)
		So(v.Value, ShouldResemble, []bool{true, false, true})
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal BoolArray", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("[true,false,true]"))
	})
}

func TestTypes_IntArray(t *testing.T) {
	var v httpin.IntArray

	Convey("Unmarshal IntArray", t, func() {
		So(json.Unmarshal([]byte("[9, 12, 1024]"), &v), ShouldBeNil)
		So(v.Value, ShouldResemble, []int{9, 12, 1024})
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal IntArray", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("[9,12,1024]"))
	})
}

func TestTypes_Float32Array(t *testing.T) {
	var v httpin.Float32Array

	Convey("Unmarshal Float32Array", t, func() {
		So(json.Unmarshal([]byte("[0.618, 2.718, 3.141]"), &v), ShouldBeNil)
		So(v.Value, ShouldResemble, []float32{0.618, 2.718, 3.141})
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal Float32Array", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("[0.618,2.718,3.141]"))
	})
}

func shouldTimeArrayEqual(actual interface{}, expected ...interface{}) string {
	actualTimes, isTimeArray := actual.([]time.Time)
	if !isTimeArray {
		return "actual is not []time.Time"
	}
	expectedTimes, isTimeArray := expected[0].([]time.Time)
	if !isTimeArray {
		return "expected is not []time.Time"
	}

	if len(actualTimes) != len(expectedTimes) {
		return fmt.Sprintf("length doesn't match, actual %d != expected %d", len(actualTimes), len(expectedTimes))
	}

	for i, actualItem := range actualTimes {
		expectedItem := expectedTimes[i]
		if !actualItem.Equal(expectedItem) {
			return fmt.Sprintf("item %d doesn't match, actual %v != expected %v", i, actualItem, expectedItem)
		}
	}

	return ""
}

func TestTypes_TimeArray(t *testing.T) {
	var v httpin.TimeArray

	Convey("Unmarshal TimeArray", t, func() {
		So(json.Unmarshal([]byte("[ \"1991-11-10T08:00:00+08:00\", \"1991-06-28T06:00:00+00:00\" ]"), &v), ShouldBeNil)
		So(v.Value, shouldTimeArrayEqual, []time.Time{
			time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("E8", 8*3600)),
			time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC),
		})
		So(v.Valid, ShouldBeFalse)
	})

	Convey("Marshal TimeArray", t, func() {
		bs, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(bs, ShouldResemble, []byte("[\"1991-11-10T08:00:00+08:00\",\"1991-06-28T06:00:00Z\"]"))
	})
}
