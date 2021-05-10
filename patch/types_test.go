package patch

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type UserPatch struct {
	Name     String `json:"name"`
	IsMember Bool   `json:"is_member"`
}

func TestTypes_DncodeInStructFields(t *testing.T) {
	Convey("Valid should be true if field exists, otherwise false", t, func() {
		var p UserPatch
		So(json.Unmarshal([]byte(`{"name": "ggicci"}`), &p), ShouldBeNil)
		So(p.Name.Valid, ShouldBeTrue)
		So(p.IsMember.Valid, ShouldBeFalse)
	})

	Convey("Unmarshal failed field also causes Valid to be false", t, func() {
		var p UserPatch
		So(json.Unmarshal([]byte(`{"name": "ggicc", "is_member": "what"}`), &p), ShouldBeError)
		So(p.Name.Valid, ShouldBeTrue)
		So(p.IsMember.Valid, ShouldBeFalse)
	})
}

func TestTypes(t *testing.T) {
	var testcases = []struct {
		Content  string
		Expected interface{}
	}{
		{"true", Bool{true, true}},
		{"2045", Int{2045, true}},
		{"127", Int8{127, true}},
		{"32767", Int16{32767, true}},
		{"2147483647", Int32{2147483647, true}},
		{"9223372036854775807", Int64{9223372036854775807, true}},
		{"2045", Uint{2045, true}},
		{"255", Uint8{255, true}},
		{"65535", Uint16{65535, true}},
		{"4294967295", Uint32{4294967295, true}},
		{"18446744073709551615", Uint64{18446744073709551615, true}},
		{"3.14", Float32{3.14, true}},
		{"3.14", Float64{3.14, true}},
		{"\"hello\"", String{"hello", true}},
		{`[true,false]`, BoolArray{[]bool{true, false}, true}},
		{"[1,2,3]", IntArray{[]int{1, 2, 3}, true}},
		{"[1,2,3]", Int8Array{[]int8{1, 2, 3}, true}},
		{"[1,2,3]", Int16Array{[]int16{1, 2, 3}, true}},
		{"[1,2,3]", Int32Array{[]int32{1, 2, 3}, true}},
		{"[1,2,3]", Int64Array{[]int64{1, 2, 3}, true}},
		{"[1,2,3]", UintArray{[]uint{1, 2, 3}, true}},
		{"[1,2,3]", Uint16Array{[]uint16{1, 2, 3}, true}},
		{"[1,2,3]", Uint32Array{[]uint32{1, 2, 3}, true}},
		{"[1,2,3]", Uint64Array{[]uint64{1, 2, 3}, true}},
		{"[0.618,1,3.14]", Float32Array{[]float32{0.618, 1, 3.14}, true}},
		{"[0.618,1,3.14]", Float64Array{[]float64{0.618, 1, 3.14}, true}},
		{`["hello","world"]`, StringArray{[]string{"hello", "world"}, true}},
	}

	Convey("Test basic types", t, func() {
		for _, testcase := range testcases {
			Convey(fmt.Sprintf("Unmarshal type: %T", testcase.Expected), func() {
				rv := reflect.New(reflect.TypeOf(testcase.Expected))
				So(json.Unmarshal([]byte(testcase.Content), rv.Interface()), ShouldBeNil)
				So(rv.Elem().Interface(), ShouldResemble, testcase.Expected)
			})

			Convey(fmt.Sprintf("Marshal type: %T", testcase.Expected), func() {
				bs, err := json.Marshal(&testcase.Expected)
				So(err, ShouldBeNil)
				So(string(bs), ShouldEqual, testcase.Content)
			})
		}
	})

	// https://golang.org/pkg/encoding/json/#Marshal
	// Array and slice values encode as JSON arrays, except that []byte encodes
	// as a base64-encoded string, and a nil slice encodes as the null JSON
	// value.
	// uint8       the set of all unsigned  8-bit integers (0 to 255)
	// byte        alias for uint8
	Convey("Type Uint8Array ([]byte, []uint8)", t, func() {
		Convey("unmarshal", func() {
			var bs Uint8Array
			So(json.Unmarshal([]byte("[1, 2, 3]"), &bs), ShouldBeNil)
			So(bs, ShouldResemble, Uint8Array{[]uint8{1, 2, 3}, true})
		})
		Convey("marshal", func() {
			var bs = Uint8Array{[]uint8{1, 2, 3}, true}
			out, err := json.Marshal(bs)
			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `"AQID"`) // base64
		})
	})

	Convey("Type Time", t, func() {
		Convey("unmarshal", func() {
			var t Time
			So(json.Unmarshal([]byte("\"1991-11-10T08:00:00+08:00\""), &t), ShouldBeNil)
			So(t.Value, ShouldEqual, time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("E8", +8*3600)))
			So(t.Valid, ShouldBeTrue)
		})
		Convey("marshal", func() {
			var t = Time{time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("E8", +8*3600)), true}
			out, err := json.Marshal(&t)
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []byte("\"1991-11-10T08:00:00+08:00\""))
		})
	})

	Convey("Type TimeArray", t, func() {
		var (
			ts = []time.Time{
				time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("E8", +8*3600)),
				time.Date(1991, 6, 28, 14, 0, 0, 0, time.FixedZone("E8", +8*3600)),
			}
			tsBytes = []byte(`["1991-11-10T08:00:00+08:00","1991-06-28T14:00:00+08:00"]`)
		)

		Convey("unmarshal", func() {
			var t TimeArray
			So(json.Unmarshal(tsBytes, &t), ShouldBeNil)
			So(t.Value, shouldTimeArrayEqual, ts)
			So(t.Valid, ShouldBeTrue)
		})
		Convey("marshal", func() {
			var t = TimeArray{ts, true}
			out, err := json.Marshal(&t)
			So(err, ShouldBeNil)
			So(out, ShouldResemble, tsBytes)
		})
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
