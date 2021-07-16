package httpin

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
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

type Pagination struct {
	Page    int `in:"form=page,page_index,index"`
	PerPage int `in:"form=per_page,page_size"`
}

type Authorization struct {
	AccessToken string `in:"form=access_token;header=x-api-token"`
}

type ProductQuery struct {
	CreatedAt time.Time `in:"form=created_at;required"`
	Color     string    `in:"form=colour,color"`
	IsSoldout bool      `in:"form=is_soldout"`
	SortBy    []string  `in:"form=sort_by"`
	SortDesc  []bool    `in:"form=sort_desc"`
	Pagination
	Authorization
}

type ObjectID struct {
	Timestamp [4]byte
	Mid       [3]byte
	Pid       [2]byte
	Counter   [3]byte
}

type Cursor struct {
	AfterMarker  ObjectID `in:"form=after"`
	BeforeMarker ObjectID `in:"form=before"`
	Limit        int      `in:"form=limit"`
}

type ThingWithInvalidDirectives struct {
	Sequence string `in:"form=seq;base58_to_integer"`
}

type ThingWithUnsupportedCustomType struct {
	Cursor
}

type ThingWithUnsupportedCustomTypeOfSliceField struct {
	IdList []ObjectID `in:"form=id[]"`
}

func TestEngine(t *testing.T) {
	Convey("New engine with non-struct type", t, func() {
		core, err := New(string("hello"))
		So(core, ShouldBeNil)
		So(errors.Is(err, ErrUnsupporetedType), ShouldBeTrue)
	})

	Convey("New engine with unregistered executor", t, func() {
		core, err := New(ThingWithInvalidDirectives{})
		So(core, ShouldBeNil)
		So(errors.Is(err, ErrUnregisteredExecutor), ShouldBeTrue)
	})

	Convey("New engine with same type should hit cache", t, func() {
		core1, err := New(ProductQuery{})
		So(err, ShouldBeNil)
		core2, err := New(ProductQuery{})
		So(err, ShouldBeNil)
		core3, err := New(&ProductQuery{})
		So(err, ShouldBeNil)
		core4, err := New(&ProductQuery{}, WithErrorStatusCode(400))
		So(err, ShouldBeNil)
		So(core1.tree, ShouldPointTo, core2.tree)
		So(core2.tree, ShouldPointTo, core3.tree)
		So(core3.tree, ShouldPointTo, core4.tree)
	})

	Convey("Very basic and normal case", t, func() {
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

	Convey("Embedded field should work", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"created_at": {"1991-11-10T08:00:00+08:00"},
			"color":      {"red"},
			"is_soldout": {"true"},
			"sort_by":    {"id", "quantity"},
			"sort_desc":  {"0", "true"},
			"page":       {"1"},
			"per_page":   {"20"},
		}
		expected := &ProductQuery{
			CreatedAt: time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC),
			Color:     "red",
			IsSoldout: true,
			SortBy:    []string{"id", "quantity"},
			SortDesc:  []bool{false, true},
			Pagination: Pagination{
				Page:    1,
				PerPage: 20,
			},
		}
		core, err := New(ProductQuery{})
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, expected)
	})

	Convey("Required field is missing", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"color":      {"red"},
			"is_soldout": {"true"},
			"sort_by":    {"id", "quantity"},
			"sort_desc":  {"0", "true"},
			"page":       {"1"},
			"per_page":   {"20"},
		}
		core, err := New(&ProductQuery{}) // struct pointer also works
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(got, ShouldBeNil)
		So(errors.Is(err, ErrMissingField), ShouldBeTrue)

		var invalidField *InvalidFieldError
		So(errors.As(err, &invalidField), ShouldBeTrue)
		So(invalidField.Source, ShouldEqual, "required")
		So(invalidField.Value, ShouldBeNil)
	})

	Convey("Non-required fields can be absent", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"created_at": {"1991-11-10T08:00:00+08:00"},
			"is_soldout": {"true"},
			"page":       {"1"},
			"per_page":   {"20"},
		}
		expected := &ProductQuery{
			CreatedAt: time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC),
			Color:     "",
			IsSoldout: true,
			Pagination: Pagination{
				Page:    1,
				PerPage: 20,
			},
		}
		core, err := New(ProductQuery{})
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, expected)
	})

	Convey("Unsupported custom type", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"uid":   {"ggicci"},
			"after": {"5cb71995ad763f7f1717c9eb"},
			"limit": {"50"},
		}
		core, err := New(ThingWithUnsupportedCustomType{})
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(got, ShouldBeNil)
		So(errors.Is(err, ErrUnsupporetedType), ShouldBeTrue)
	})

	Convey("Unsupported custom type of slice field", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"id[]": {
				"5cb71995ad763f7f1717c9eb",
				"60922dd8940cf19c30bba50c",
				"6093a70fdb597d966944c125",
			},
		}
		core, err := New(ThingWithUnsupportedCustomTypeOfSliceField{})
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(got, ShouldBeNil)
		So(errors.Is(err, ErrUnsupporetedType), ShouldBeTrue)
	})

	Convey("Meet invalid value for a key", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"created_at": {"1991-11-10T08:00:00+08:00"},
			"is_soldout": {"zero"}, // invalid
		}
		core, err := New(ProductQuery{})
		So(err, ShouldBeNil)
		_, err = core.Decode(r)
		So(err, ShouldBeError)
		var invalidField *InvalidFieldError
		So(errors.As(err, &invalidField), ShouldBeTrue)
		So(invalidField.Field, ShouldEqual, "IsSoldout")
		So(invalidField.Source, ShouldEqual, "form")
		So(invalidField.Value, ShouldEqual, "zero")
	})

	Convey("Meet invalid values for a key (of slice type)", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"created_at": {"1991-11-10T08:00:00+08:00"},
			"sort_desc":  {"true", "zero", "0"}, // invalid value "zero"
		}
		core, err := New(ProductQuery{})
		So(err, ShouldBeNil)
		_, err = core.Decode(r)
		var invalidField *InvalidFieldError
		So(errors.As(err, &invalidField), ShouldBeTrue)
		So(invalidField.Field, ShouldEqual, "SortDesc")
		So(invalidField.Source, ShouldEqual, "form")
		So(invalidField.Value, ShouldResemble, []string{"true", "zero", "0"})
		So(err.Error(), ShouldContainSubstring, "at index 1")
	})

	Convey("Custom decoder should work", t, func() {
		var boolType = reflect.TypeOf(bool(true))
		RegisterTypeDecoder(boolType, TypeDecoderFunc(DecodeCustomBool))
		type BoolInput struct {
			IsMember bool `in:"form=is_member"`
		}
		core, _ := New(BoolInput{})
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{"is_member": {"yes"}}
		got, err := core.Decode(r)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, &BoolInput{IsMember: true})
		delete(decoders, boolType) // remove the custom decoder
	})
}
