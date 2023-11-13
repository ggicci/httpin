package core

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/owl"
	"github.com/stretchr/testify/assert"
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

func TestNew_WithNonStructType(t *testing.T) {
	co, err := New(string("hello"))
	assert.Nil(t, co)
	assert.ErrorIs(t, err, owl.ErrUnsupportedType)
}

func TestNew_WithUnregisteredExecutor(t *testing.T) {
	type ThingWithInvalidDirectives struct {
		Sequence string `in:"form=seq;base58_to_integer"`
	}

	co, err := New(ThingWithInvalidDirectives{})
	assert.Nil(t, co)
	assert.ErrorContains(t, err, "unregistered directive")
	assert.ErrorContains(t, err, "base58_to_integer")
}

func CustomErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	var invalidFieldError *InvalidFieldError
	if errors.As(err, &invalidFieldError) {
		rw.WriteHeader(http.StatusBadRequest) // status: 400
		io.WriteString(rw, invalidFieldError.Error())
		return
	}
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // status: 500
}

func TestNew_hitCachedResolverOfSameInputType(t *testing.T) {
	assert := assert.New(t)

	type Query struct{}
	core1, err := New(Query{})
	assert.NoError(err)

	core2, err := New(Query{})
	assert.NoError(err)

	assert.Equal(core1.resolver, core2.resolver)

	core3, err := New(&Query{}, WithErrorHandler(CustomErrorHandler))
	assert.NoError(err)
	assert.Equal(core1.resolver, core3.resolver)
}

func TestCore_Decode_EmbeddedStruct(t *testing.T) {
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
	co, err := New(ProductQuery{})
	assert.NoError(t, err)
	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.(*ProductQuery))
}

func TestCore_Decode_DecodeError_InvalidSingleValue(t *testing.T) {
	assert := assert.New(t)
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"created_at": {"1991-11-10T08:00:00+08:00"},
		"is_soldout": {"zero"}, // invalid
	}
	co, err := New(ProductQuery{})
	assert.NoError(err)
	_, err = co.Decode(r)
	var invalidField *InvalidFieldError
	assert.ErrorAs(err, &invalidField)
	assert.Equal("IsSoldout", invalidField.Field)
	assert.Equal("form", invalidField.Source)
	assert.Equal("is_soldout", invalidField.Key)
	assert.Equal([]string{"zero"}, invalidField.Value)
}

func TestCore_Decode_DecodeError_InvalidValueInSlice(t *testing.T) {
	assert := assert.New(t)
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"created_at": {"1991-11-10T08:00:00+08:00"},
		"sort_desc":  {"true", "zero", "0"}, // invalid value "zero"
	}
	co, err := New(ProductQuery{})
	assert.NoError(err)
	_, err = co.Decode(r)
	var invalidField *InvalidFieldError
	assert.ErrorAs(err, &invalidField)
	assert.Equal("SortDesc", invalidField.Field)
	assert.Equal("form", invalidField.Source)
	assert.Equal("sort_desc", invalidField.Key)
	assert.Equal([]string{"true", "zero", "0"}, invalidField.Value)
	assert.ErrorContains(err, "at index 1")
}

func TestCore_Decode_ErrUnsupporetedType(t *testing.T) {
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

	type Payload struct {
		IdList []ObjectID `in:"form=id[]"`
	}

	// Unsupported custom type as field.
	func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"uid":   {"ggicci"},
			"after": {"5cb71995ad763f7f1717c9eb"},
			"limit": {"50"},
		}
		co, err := New(Cursor{})
		assert.NoError(t, err)
		got, err := co.Decode(r)
		assert.ErrorIs(t, err, ErrUnsupportedType)
		assert.ErrorContains(t, err, "ObjectID")
		assert.Nil(t, got)
	}()

	// Slice of unsupported type.
	func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"id[]": {
				"5cb71995ad763f7f1717c9eb",
				"60922dd8940cf19c30bba50c",
				"6093a70fdb597d966944c125",
			},
		}
		co, err := New(Payload{})
		assert.NoError(t, err)
		got, err := co.Decode(r)
		assert.ErrorIs(t, err, ErrUnsupportedType)
		assert.ErrorContains(t, err, "ObjectID")
		assert.Nil(t, got)
	}()
}

func TestCore_Decode_UnexportedFields(t *testing.T) {
	type ThingWithUnexportedFields struct {
		Name    string `in:"form=name"`
		display string `in:"form=display"` // unexported field
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"name":    []string{"ggicci"},
		"display": []string{"Ggicci T'ang"},
	}
	expected := &ThingWithUnexportedFields{
		Name:    "ggicci",
		display: "",
	}
	co, err := New(ThingWithUnexportedFields{})
	assert.NoError(t, err)
	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.(*ThingWithUnexportedFields))
}

func TestCore_Decode_CustomDecoder_TypeDecoder(t *testing.T) {
	RegisterDecoder[bool](myBoolDecoder) // usually done in init()
	RegisterDecoder[Place](myPlaceDecoder)

	type Input struct {
		IsMember           bool   `in:"form=is_member"`
		RegisterationPlace *Place `in:"form=registration_place"`
	}
	co, err := New(Input{})
	assert.NoError(t, err)

	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"is_member":          {"yes"},
		"registration_place": {"Canada.Toronto"},
	}

	expected := &Input{
		IsMember:           true,
		RegisterationPlace: &Place{Country: "Canada", City: "Toronto"},
	}

	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	removeTypeDecoder[bool]() // remove the custom decoder
	removeTypeDecoder[Place]()
}

func TestCore_Decode_CustomDecoder_RegisterThePointerType(t *testing.T) {
	RegisterDecoder[*Place](myPlacePointerDecoder) // usually done in init()
	type Input struct {
		BornPlace *Place `in:"form=born_place"`
		LivePlace Place  `in:"form=live_place"`
	}
	co, err := New(Input{})
	assert.NoError(t, err)

	r := newMultipartFormRequestFromMap(map[string]any{
		"born_place": "China.Huzhou",
		"live_place": "Canada.Toronto",
	})
	expected := &Input{
		BornPlace: &Place{Country: "China", City: "Huzhou"},
		LivePlace: Place{Country: "Canada", City: "Toronto"},
	}

	gotValue, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, gotValue)

	removeTypeDecoder[*Place]() // remove the custom decoder
	removeTypeDecoder[Place]()
}

func TestCore_Decode_PointerTypes(t *testing.T) {
	assert := assert.New(t)

	RegisterDecoder[Place](myPlaceDecoder) // usually done in init()
	type Input struct {
		IsMember       *bool  `in:"form=is_member"`
		Limit          *int   `in:"form=limit"`
		LastAccessFrom *Place `in:"form=_laf"`
	}
	co, err := New(Input{})
	assert.NoError(err)

	// Missing fields.
	r := newMultipartFormRequestFromMap(map[string]any{
		"is_member": "true",
	})
	gotValue, err := co.Decode(r)
	assert.NoError(err)
	got := gotValue.(*Input)
	assert.Equal(true, *got.IsMember)
	assert.Nil(got.Limit)
	assert.Nil(got.LastAccessFrom)

	// All fields.
	r = newMultipartFormRequestFromMap(map[string]any{
		"is_member": "true",
		"limit":     "10",
		"_laf":      "Canada.Toronto",
	})
	gotValue, err = co.Decode(r)
	assert.NoError(err)
	got = gotValue.(*Input)
	assert.Equal(true, *got.IsMember)
	assert.Equal(10, *got.Limit)
	assert.Equal(Place{Country: "Canada", City: "Toronto"}, *got.LastAccessFrom)

	// Invalid value.
	r = newMultipartFormRequestFromMap(map[string]any{
		"_laf": "Canada", // invalid value
	})
	gotValue, err = co.Decode(r)
	assert.Nil(gotValue)
	var ife *InvalidFieldError
	assert.ErrorAs(err, &ife)
	assert.Equal("_laf", ife.Key)
	assert.Equal([]string{"Canada"}, ife.Value)
	assert.Equal("form", ife.Source)
	assert.ErrorContains(err, "invalid place")

	removeTypeDecoder[Place]()
}

type CustomNamedDecoderInput struct {
	Name             string      `in:"form=name"`
	Birthday         time.Time   `in:"form=birthday;decoder=decodeMyDate"` // use named decoder "decodeMyDate"
	EffectiveBetween []time.Time `in:"form=effective_between;decoder=decodeMyDate"`
	CreatedBetween   []time.Time `in:"form=created_between"`
}

func TestCore_Decode_CustomDecoder_NamedDecoder(t *testing.T) {
	RegisterNamedDecoder[time.Time]("decodeMyDate", myDateDecoder, true) // usually done in init()

	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"name":              {"Ggicci"},
		"birthday":          {"1991-11-10"},
		"effective_between": {"2021-04-12", "2025-04-12"},
		"created_between":   {"2021-01-01T08:00:00+08:00", "2022-01-01T08:00:00+08:00"},
	}
	expected := &CustomNamedDecoderInput{
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
	}

	co, err := New(CustomNamedDecoderInput{})
	assert.NoError(t, err)
	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.(*CustomNamedDecoderInput))
}

func TestCore_Decode_CustomDecoder_NamedDecoder_ErrMissingDecoderName(t *testing.T) {
	type GenderType string
	type Input struct {
		Gender GenderType `in:"form=gender;decoder"`
	}

	co, err := New(Input{})
	assert.ErrorContains(t, err, "missing decoder name")
	assert.Nil(t, co)
}

func TestCore_Decode_CustomDecoder_NamedDecoder_ErrUnregisteredDecoder(t *testing.T) {
	type GenderType string
	type Input struct {
		Gender GenderType `in:"form=gender;decoder=notfound"` // cause ErrUnregisteredDecoder
	}
	co, err := New(Input{})
	assert.ErrorContains(t, err, "unregistered decoder: \"notfound\"")
	assert.Nil(t, co)
}

func TestCore_Decode_CustomDecoder_NamedDecoder_ErrCannotSpecifyOnFileTypeFields(t *testing.T) {
	RegisterNamedDecoder[time.Time]("decodeMyDate", myDateDecoder, true) // usually done in init()
	type FunnyFile struct{}
	defaultRegistry.fileTypes[internal.TypeOf[*FunnyFile]()] = nil // fake a registered file type
	type Input struct {
		Avatar *FunnyFile `in:"form=avatar;decoder=decodeMyDate"`
	}
	co, err := New(Input{})
	assert.ErrorContains(t, err, "cannot use decoder directive on a file type field")
	assert.Nil(t, co)
}

func TestCore_Decode_CustomDecoder_NamedDecoder_ErrValueTypeMismatch(t *testing.T) {
	RegisterNamedDecoder[time.Time]("decodeMyDate", myDateDecoder, true) // usually done in init()

	type Input struct {
		Birthday string `in:"form=birthday;decoder=decodeMyDate"` // cause ErrValueTypeMismatch
	}

	co, err := New(Input{})
	assert.NoError(t, err)
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{"birthday": {"1991-11-10"}} // birthday is string, not time.Time
	_, err = co.Decode(r)
	assert.ErrorIs(t, err, ErrTypeMismatch)
	assert.ErrorContains(t, err, "birthday")
	assert.ErrorContains(t, err, "string")
	assert.ErrorContains(t, err, "time.Time")
}

func TestCore_Decode_CustomDecoder_NamedDecoder_DecodeError(t *testing.T) {
	RegisterNamedDecoder[time.Time]("decodeMyDate", myDateDecoder, true) // usually done in init()

	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"name":     {"Ggicci"},
		"birthday": {"1991-11-10 08:00:00"}, // invalid date format
	}

	co, err := New(CustomNamedDecoderInput{})
	assert.NoError(t, err)

	got, err := co.Decode(r)
	var invalidDate *InvalidDate
	assert.ErrorAs(t, err, &invalidDate)
	assert.ErrorContains(t, err, "invalid date: \"1991-11-10 08:00:00\"")
	assert.Nil(t, got)
}
