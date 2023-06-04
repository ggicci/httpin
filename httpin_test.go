package httpin

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

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

type ThingWithUnexportedFields struct {
	Name    string `in:"form=name"`
	display string // unexported field
}

func decodeMyDate(value string) (interface{}, error) {
	return time.Parse("2006-01-02", value)
}

type ThingWithNamedDecoder struct {
	Name             string      `in:"form=name"`
	Birthday         time.Time   `in:"form=birthday;decoder=decodeMyDate"`
	EffectiveBetween []time.Time `in:"form=effective_between;decoder=decodeMyDate"`
	CreatedBetween   []time.Time `in:"form=created_between"`
}

type InvalidName struct {
	Name string
}

func (e *InvalidName) Error() string {
	return fmt.Sprintf("name '%s' is invalid", e.Name)
}

func decodeName(value string) (interface{}, error) {
	return nil, &InvalidName{Name: value}
}

type ThingWithInvalidNamedDecoder struct {
	Name string `in:"form=name;decoder=decodeName"`
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
		core4, err := New(&ProductQuery{}, WithErrorHandler(CustomErrorHandler))
		So(err, ShouldBeNil)
		So(core1.tree, ShouldEqual, core2.tree)
		So(core2.tree, ShouldEqual, core3.tree)
		So(core3.tree, ShouldEqual, core4.tree)
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

	Convey("Unexported fields should be ignored", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"name": []string{"ggicci"},
		}
		expected := &ThingWithUnexportedFields{
			Name:    "ggicci",
			display: "",
		}
		core, err := New(ThingWithUnexportedFields{})
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
		RegisterTypeDecoder(boolType, ValueTypeDecoderFunc(decodeCustomBool))
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

	Convey("Can specify a named decoder for a field rather than being auto-selected by field type", t, func() {
		ReplaceNamedDecoder("decodeMyDate", ValueTypeDecoderFunc(decodeMyDate))

		engine, err := New(ThingWithNamedDecoder{})
		So(err, ShouldBeNil)

		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"name":              {"Ggicci"},
			"birthday":          {"1991-11-10"},
			"effective_between": {"2021-04-12", "2025-04-12"},
			"created_between":   {"2021-01-01T08:00:00+08:00", "2022-01-01T08:00:00+08:00"},
		}
		got, err := engine.Decode(r)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, &ThingWithNamedDecoder{
			Name:     "Ggicci",
			Birthday: time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC),
			EffectiveBetween: []time.Time{
				time.Date(2021, 4, 12, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 4, 12, 0, 0, 0, 0, time.UTC),
			},
			CreatedBetween: []time.Time{
				time.Date(2021, 1, 1, 8, 0, 0, 0, time.FixedZone("E8", +8*3600)).UTC(),
				time.Date(2022, 1, 1, 8, 0, 0, 0, time.FixedZone("E8", +8*3600)).UTC(),
			},
		})
	})

	Convey("Can unwrap custom errors from named decoder", t, func() {
		ReplaceNamedDecoder("decodeName", ValueTypeDecoderFunc(decodeName))

		engine, err := New(ThingWithInvalidNamedDecoder{})
		So(err, ShouldBeNil)

		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"name": {"Dragomeat"},
		}
		_, err = engine.Decode(r)
		So(err, ShouldBeError)
		var invalidName *InvalidName
		So(errors.As(err, &invalidName), ShouldBeTrue)
		So(invalidName.Name, ShouldEqual, "Dragomeat")
	})
}

func TestDecode(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"page":     {"1"},
		"per_page": {"100"},
	}
	expected := &Pagination{
		Page:    1,
		PerPage: 100,
	}

	Convey("Decode passing a pointer to a pointer of struct instance", t, func() {
		input := &Pagination{}
		So(Decode(r, &input), ShouldBeNil)
		So(input, ShouldResemble, expected)
	})

	Convey("Decode passing a pointer to a struct instance", t, func() {
		input := Pagination{}
		So(Decode(r, &input), ShouldBeNil)
		So(input, ShouldResemble, *expected)
	})

	Convey("Decode passing a struct instance should fail", t, func() {
		input := Pagination{}
		So(Decode(r, input), ShouldBeError)
	})
}
