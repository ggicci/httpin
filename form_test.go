package httpin

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// ChaosQuery is designed to make the normal case test coverage higher.
type ChaosQuery struct {
	// Basic Types
	BoolValue       bool       `in:"form=bool"`
	IntValue        int        `in:"form=int"`
	Int8Value       int8       `in:"form=int8"`
	Int16Value      int16      `in:"form=int16"`
	Int32Value      int32      `in:"form=int32"`
	Int64Value      int64      `in:"form=int64"`
	UintValue       uint       `in:"form=uint"`
	Uint8Value      uint8      `in:"form=uint8"`
	Uint16Value     uint16     `in:"form=uint16"`
	Uint32Value     uint32     `in:"form=uint32"`
	Uint64Value     uint64     `in:"form=uint64"`
	Float32Value    float32    `in:"form=float32"`
	Float64Value    float64    `in:"form=float64"`
	Complex64Value  complex64  `in:"form=complex64"`
	Complex128Value complex128 `in:"form=complex128"`
	StringValue     string     `in:"form=string"`

	// Time Type
	TimeValue time.Time `in:"form=time"`

	// Array
	BoolList   []bool      `in:"form=bools"`
	IntList    []int       `in:"form=ints"`
	FloatList  []float64   `in:"form=floats"`
	StringList []string    `in:"form=strings"`
	TimeList   []time.Time `in:"form=times"`
}

func TestDirectiveForm(t *testing.T) {
	Convey("Very basic and normal cases", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"bool":       {"true"},
			"int":        {"9"},
			"int8":       {"14"},
			"int16":      {"841"},
			"int32":      {"193"},
			"int64":      {"475"},
			"uint":       {"11"},
			"uint8":      {"4"},
			"uint16":     {"48"},
			"uint32":     {"9583"},
			"uint64":     {"183471"},
			"float32":    {"3.14"},
			"float64":    {"0.618"},
			"complex64":  {"1+4i"},
			"complex128": {"-6+17i"},
			"string":     {"doggy"},
			"time":       {"1991-11-10T08:00:00+08:00"},
			"bools":      {"true", "false", "0", "1"},
			"ints":       {"9", "9", "6"},
			"floats":     {"0", "0.5", "1"},
			"strings":    {"Life", "is", "a", "Miracle"},
			"times":      {"2000-01-02T15:04:05-07:00", "678088800"},
		}
		expected := &ChaosQuery{
			BoolValue:       true,
			IntValue:        9,
			Int8Value:       14,
			Int16Value:      841,
			Int32Value:      193,
			Int64Value:      475,
			UintValue:       11,
			Uint8Value:      4,
			Uint16Value:     48,
			Uint32Value:     9583,
			Uint64Value:     183471,
			Float32Value:    3.14,
			Float64Value:    0.618,
			Complex64Value:  1 + 4i,
			Complex128Value: -6 + 17i,
			StringValue:     "doggy",
			TimeValue:       time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC),
			BoolList:        []bool{true, false, false, true},
			IntList:         []int{9, 9, 6},
			FloatList:       []float64{0.0, 0.5, 1.0},
			StringList:      []string{"Life", "is", "a", "Miracle"},
			TimeList: []time.Time{
				time.Date(2000, 1, 2, 22, 4, 5, 0, time.UTC),
				time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC),
			},
		}

		core, err := New(ChaosQuery{})
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, expected)
	})
}
