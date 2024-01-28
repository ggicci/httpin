package core

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ggicci/httpin/internal"
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

// Place is a custom type that implements Stringable interface.
type Place struct {
	Country string
	City    string
}

func (p Place) ToString() (string, error) {
	return fmt.Sprintf("%s.%s", p.Country, p.City), nil
}

func (p *Place) FromString(value string) error {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid place: %q", value)
	}
	*p = Place{Country: parts[0], City: parts[1]}
	return nil
}

func TestNew_WithNonStructType(t *testing.T) {
	co, err := New(string("hello"))
	assert.Nil(t, co)
	assert.ErrorIs(t, err, ErrUnsupportedType)
}

func TestNew_ErrUnregisteredDirective(t *testing.T) {
	type ThingWithInvalidDirectives struct {
		Sequence string `in:"form=seq;base58_to_integer"`
	}

	co, err := New(ThingWithInvalidDirectives{})
	assert.Nil(t, co)
	assert.ErrorContains(t, err, "unregistered directive")
	assert.ErrorContains(t, err, "base58_to_integer")
}

func TestNew_WithNamedCoder_ErrMissingCoderName(t *testing.T) {
	type Input struct {
		Gender string `in:"form=gender;decoder"`
	}

	co, err := New(Input{})
	assert.ErrorContains(t, err, "directive decoder: missing coder name")
	assert.Nil(t, co)
}

func TestNew_WithNamedCoder_ErrUnregisteredCoder(t *testing.T) {
	type Input struct {
		Gender string `in:"form=gender;coder=gender"`
	}
	co, err := New(Input{})
	assert.ErrorContains(t, err, "directive coder: unregistered coder: \"gender\"")
	assert.Nil(t, co)
}

func TestNew_WithNamedCoder_ErrCannotSpecifyOnFileTypeFields(t *testing.T) {
	registerMyDate()
	type FunnyFile struct{}
	fileTypes[internal.TypeOf[*FunnyFile]()] = struct{}{} // fake a registered file type
	type Input struct {
		Avatar *FunnyFile `in:"form=avatar;coder=mydate"`
	}
	co, err := New(Input{})
	assert.ErrorContains(t, err, "directive coder: cannot be used on a file type field")
	assert.Nil(t, co)
	removeFileType[*FunnyFile]()
	unregisterMyDate()
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

// We only build the resolver once for each input type.
func TestNew_HitCachedResolverOfSameInputType(t *testing.T) {
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

func TestCore_Decode_ErrInvalidSingleValue(t *testing.T) {
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
	assert.Equal("form", invalidField.Directive)
	assert.Equal("is_soldout", invalidField.Key)
	assert.Equal([]string{"zero"}, invalidField.Value)
}

func TestCore_Decode_ErrInvalidValueInSlice(t *testing.T) {
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
	assert.Equal("form", invalidField.Directive)
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

func TestCore_Decode_SkipUnexportedFields(t *testing.T) {
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

func TestCore_Decode_PointerTypes(t *testing.T) {
	assert := assert.New(t)

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
	assert.Equal("form", ife.Directive)
	assert.ErrorContains(err, "invalid place")
}

type CommaSeparatedIntegerArray struct {
	Value []int
}

func (a CommaSeparatedIntegerArray) ToString() (string, error) {
	var res = make([]string, len(a.Value))
	for i := range a.Value {
		res[i] = strconv.Itoa(a.Value[i])
	}
	return strings.Join(res, ","), nil
}

func (pa *CommaSeparatedIntegerArray) FromString(value string) error {
	a := CommaSeparatedIntegerArray{}
	values := strings.Split(value, ",")
	a.Value = make([]int, len(values))
	for i := range values {
		if value, err := strconv.Atoi(values[i]); err != nil {
			return err
		} else {
			a.Value[i] = value
		}
	}
	*pa = a
	return nil
}

func TestCore_Decode_CustomTypeSliceValueWrapper(t *testing.T) {
	assert := assert.New(t)

	type Input struct {
		Ids CommaSeparatedIntegerArray `in:"form=ids"`
	}
	co, err := New(Input{})
	assert.NoError(err)

	// Missing fields.
	r := newMultipartFormRequestFromMap(map[string]any{
		"ids": "1,2,3",
	})
	gotValue, err := co.Decode(r)
	assert.NoError(err)
	got := gotValue.(*Input)
	assert.Equal([]int{1, 2, 3}, got.Ids.Value)
}

// Test: register named coders and use them in the "coder" directive,
// i.e. customizing the encoding/decoding for a specific struct field.

type NamedCoderInput struct {
	Name             string      `in:"form=name"`
	Birthday         time.Time   `in:"form=birthday;coder=mydate"` // use named coder "mydate"
	EffectiveBetween []time.Time `in:"form=effective_between;coder=mydate"`
	CreatedBetween   []time.Time `in:"form=created_between"`
}

func TestCore_NamedCoder(t *testing.T) {
	registerMyDate()
	co, err := New(NamedCoderInput{})
	assert.NoError(t, err)

	sampleInput := &NamedCoderInput{
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

	// Decode
	func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"name":              {"Ggicci"},
			"birthday":          {"1991-11-10"},
			"effective_between": {"2021-04-12", "2025-04-12"},
			"created_between":   {"2021-01-01T08:00:00+08:00", "2022-01-01T08:00:00+08:00"},
		}

		got, err := co.Decode(r)
		assert.NoError(t, err)
		assert.Equal(t, sampleInput, got.(*NamedCoderInput))
	}()

	// Encode / NewRequest
	func() {
		req, err := co.NewRequest("PUT", "/users/ggicci", sampleInput)
		assert.NoError(t, err)

		expected, _ := http.NewRequest("PUT", "/users/ggicci", nil)
		expectedForm := url.Values{
			"name":              {"Ggicci"},
			"birthday":          {"1991-11-10"},
			"effective_between": {"2021-04-12", "2025-04-12"},
			"created_between":   {"2021-01-01T00:00:00Z", "2022-01-01T00:00:00Z"},
		}
		expected.Body = io.NopCloser(strings.NewReader(expectedForm.Encode()))
		expected.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		assert.Equal(t, expected, req)
	}()

	unregisterMyDate()
}

func TestCore_NamedCoder_ErrTypeMismatch(t *testing.T) {
	registerMyDate()
	type Input struct {
		Birthday string `in:"form=birthday;coder=mydate"` // mydate is for time.Time, not string
	}

	co, err := New(Input{})
	assert.NoError(t, err)

	// Decode
	func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{"birthday": {"1991-11-10"}}
		_, err = co.Decode(r)
		assert.ErrorIs(t, err, internal.ErrTypeMismatch)
		assert.ErrorContains(t, err, "Birthday")
		assert.ErrorContains(t, err, "string")
		assert.ErrorContains(t, err, "time.Time")
	}()

	// Encode / NewRequest
	func() {
		payload := &Input{Birthday: "1991-11-10"}
		_, err := co.NewRequest("GET", "/", payload)
		assert.ErrorIs(t, err, internal.ErrTypeMismatch)
		assert.ErrorContains(t, err, "Birthday")
		assert.ErrorContains(t, err, "string")
		assert.ErrorContains(t, err, "time.Time")
	}()

	unregisterMyDate()
}

func TestCore_NamedCoder_DecoderError(t *testing.T) {
	registerMyDate()
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"name":     {"Ggicci"},
		"birthday": {"1991-11-10 08:00:00"}, // invalid date format
	}

	co, err := New(NamedCoderInput{})
	assert.NoError(t, err)

	got, err := co.Decode(r)
	var invalidDate *InvalidDate
	assert.ErrorAs(t, err, &invalidDate)
	assert.ErrorContains(t, err, "invalid date: \"1991-11-10 08:00:00\"")
	assert.Nil(t, got)
	unregisterMyDate()
}

func TestCore_NewRequest_NamedCoder(t *testing.T) {
	registerMyDate()

	unregisterMyDate()
}

// Test: custom type: override the default types

type YesNo bool

func (yn YesNo) ToString() (string, error) {
	if yn == true {
		return "yes", nil
	}
	return "no", nil
}

func (yn *YesNo) FromString(s string) error {
	switch s {
	case "yes":
		*yn = true
	case "no":
		*yn = false
	default:
		return fmt.Errorf("invalid YesNo value: %q", s)
	}
	return nil
}

func TestRegisterCoder_CustomType_OverrideDefaultTypeCoder(t *testing.T) {
	RegisterCoder[bool](func(b *bool) (internal.Stringable, error) {
		return (*YesNo)(b), nil
	})

	type Input struct {
		IsMember           bool   `in:"form=is_member"`
		RegisterationPlace *Place `in:"form=registration_place"`
	}
	co, err := New(Input{})
	assert.NoError(t, err)

	// Decode
	func() {
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
	}()

	// Encode / NewRequest
	func() {
		payload := &Input{
			IsMember:           true,
			RegisterationPlace: &Place{Country: "US", City: "New_York"},
		}
		expected, _ := http.NewRequest("GET", "/search", nil)
		expectedForm := url.Values{
			"is_member":          {"yes"},
			"registration_place": {"US.New_York"},
		}
		expected.Body = io.NopCloser(strings.NewReader(expectedForm.Encode()))
		expected.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		req, err := co.NewRequest("GET", "/search", payload)
		assert.NoError(t, err)
		assert.Equal(t, expected, req)
	}()

	removeType[bool]()
}

func removeType[T any]() {
	delete(customStringableAdaptors, internal.TypeOf[T]())
}

func removeNamedType(name string) {
	delete(namedStringableAdaptors, name)
}

func registerMyDate() {
	RegisterNamedCoder[time.Time]("mydate", func(t *time.Time) (Stringable, error) {
		return (*MyDate)(t), nil
	})
}

func unregisterMyDate() { removeNamedType("mydate") }
