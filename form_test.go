package httpin

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	TimeValue       time.Time  `in:"form=time"` // time type is special

	// Pointer Types
	BoolPointer       *bool       `in:"form=bool_pointer"`
	IntPointer        *int        `in:"form=int_pointer"`
	Int8Pointer       *int8       `in:"form=int8_pointer"`
	Int16Pointer      *int16      `in:"form=int16_pointer"`
	Int32Pointer      *int32      `in:"form=int32_pointer"`
	Int64Pointer      *int64      `in:"form=int64_pointer"`
	UintPointer       *uint       `in:"form=uint_pointer"`
	Uint8Pointer      *uint8      `in:"form=uint8_pointer"`
	Uint16Pointer     *uint16     `in:"form=uint16_pointer"`
	Uint32Pointer     *uint32     `in:"form=uint32_pointer"`
	Uint64Pointer     *uint64     `in:"form=uint64_pointer"`
	Float32Pointer    *float32    `in:"form=float32_pointer"`
	Float64Pointer    *float64    `in:"form=float64_pointer"`
	Complex64Pointer  *complex64  `in:"form=complex64_pointer"`
	Complex128Pointer *complex128 `in:"form=complex128_pointer"`
	StringPointer     *string     `in:"form=string_pointer"`
	TimePointer       *time.Time  `in:"form=time_pointer"`

	// Array
	BoolList   []bool      `in:"form=bools"`
	IntList    []int       `in:"form=ints"`
	FloatList  []float64   `in:"form=floats"`
	StringList []string    `in:"form=strings"`
	TimeList   []time.Time `in:"form=times"`
}

var (
	sampleChaosQuery = &ChaosQuery{
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

		BoolPointer:       asPointerValue[bool](true),
		IntPointer:        asPointerValue[int](9),
		Int8Pointer:       asPointerValue[int8](14),
		Int16Pointer:      asPointerValue[int16](841),
		Int32Pointer:      asPointerValue[int32](193),
		Int64Pointer:      asPointerValue[int64](475),
		UintPointer:       asPointerValue[uint](11),
		Uint8Pointer:      asPointerValue[uint8](4),
		Uint16Pointer:     asPointerValue[uint16](48),
		Uint32Pointer:     asPointerValue[uint32](9583),
		Uint64Pointer:     asPointerValue[uint64](183471),
		Float32Pointer:    asPointerValue[float32](3.14),
		Float64Pointer:    asPointerValue[float64](0.618),
		Complex64Pointer:  asPointerValue[complex64](1 + 4i),
		Complex128Pointer: asPointerValue[complex128](-6 + 17i),
		StringPointer:     asPointerValue[string]("doggy"),
		TimePointer:       asPointerValue[time.Time](time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC)),

		BoolList:   []bool{true, false, false, true},
		IntList:    []int{9, 9, 6},
		FloatList:  []float64{0.0, 0.5, 1.0},
		StringList: []string{"Life", "is", "a", "Miracle"},
		TimeList: []time.Time{
			time.Date(2000, 1, 2, 22, 4, 5, 0, time.UTC),
			time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC),
		},
	}
)

func TestDirectiveForm(t *testing.T) {
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

		"bool_pointer":       {"true"},
		"int_pointer":        {"9"},
		"int8_pointer":       {"14"},
		"int16_pointer":      {"841"},
		"int32_pointer":      {"193"},
		"int64_pointer":      {"475"},
		"uint_pointer":       {"11"},
		"uint8_pointer":      {"4"},
		"uint16_pointer":     {"48"},
		"uint32_pointer":     {"9583"},
		"uint64_pointer":     {"183471"},
		"float32_pointer":    {"3.14"},
		"float64_pointer":    {"0.618"},
		"complex64_pointer":  {"1+4i"},
		"complex128_pointer": {"-6+17i"},
		"string_pointer":     {"doggy"},
		"time_pointer":       {"1991-11-10T08:00:00+08:00"},

		"bools":   {"true", "false", "0", "1"},
		"ints":    {"9", "9", "6"},
		"floats":  {"0", "0.5", "1"},
		"strings": {"Life", "is", "a", "Miracle"},
		"times":   {"2000-01-02T15:04:05-07:00", "678088800"},
	}
	expected := sampleChaosQuery
	core, err := New(ChaosQuery{})
	assert.NoError(t, err)
	got, err := core.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.(*ChaosQuery))
}

func TestDirectiveForm_Encode(t *testing.T) {
	core, err := New(ChaosQuery{})
	assert.NoError(t, err)
	req, err := core.NewRequest("POST", "/signup", sampleChaosQuery)
	assert.NoError(t, err)

	expected, _ := http.NewRequest("POST", "/signup", nil)
	expected.Form = url.Values{
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
		"complex64":  {"(1+4i)"},
		"complex128": {"(-6+17i)"},
		"string":     {"doggy"},
		"time":       {"1991-11-10T00:00:00Z"},

		"bool_pointer":       {"true"},
		"int_pointer":        {"9"},
		"int8_pointer":       {"14"},
		"int16_pointer":      {"841"},
		"int32_pointer":      {"193"},
		"int64_pointer":      {"475"},
		"uint_pointer":       {"11"},
		"uint8_pointer":      {"4"},
		"uint16_pointer":     {"48"},
		"uint32_pointer":     {"9583"},
		"uint64_pointer":     {"183471"},
		"float32_pointer":    {"3.14"},
		"float64_pointer":    {"0.618"},
		"complex64_pointer":  {"(1+4i)"},
		"complex128_pointer": {"(-6+17i)"},
		"string_pointer":     {"doggy"},
		"time_pointer":       {"1991-11-10T00:00:00Z"},

		"bools":   {"true", "false", "false", "true"},
		"ints":    {"9", "9", "6"},
		"floats":  {"0", "0.5", "1"},
		"strings": {"Life", "is", "a", "Miracle"},
		"times":   {"2000-01-02T22:04:05Z", "1991-06-28T06:00:00Z"},
	}
	expected.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.Equal(t, expected, req)
}
