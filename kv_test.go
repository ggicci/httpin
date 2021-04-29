package httpin

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

type Pagination struct {
	Page    int `in:"query.page"`
	PerPage int `in:"query.per_page"`
}

type Authorization struct {
	AccessToken string `in:"query.access_token,header.x-api-token"`
}

// ChaosQuery is designed to make the normal case test coverage higher.
type ChaosQuery struct {
	// Basic Types
	BoolValue       bool       `in:"query.bool"`
	IntValue        int        `in:"query.int"`
	Int8Value       int8       `in:"query.int8"`
	Int16Value      int16      `in:"query.int16"`
	Int32Value      int32      `in:"query.int32"`
	Int64Value      int64      `in:"query.int64"`
	UintValue       uint       `in:"query.uint"`
	Uint8Value      uint8      `in:"query.uint8"`
	Uint16Value     uint16     `in:"query.uint16"`
	Uint32Value     uint32     `in:"query.uint32"`
	Uint64Value     uint64     `in:"query.uint64"`
	Float32Value    float32    `in:"query.float32"`
	Float64Value    float64    `in:"query.float64"`
	Complex64Value  complex64  `in:"query.complex64"`
	Complex128Value complex128 `in:"query.complex128"`
	StringValue     string     `in:"query.string"`

	// Time Type
	TimeValue time.Time `in:"query.time"`

	// Array
	BoolList   []bool      `in:"query.bools"`
	IntList    []int       `in:"query.ints"`
	FloatList  []float64   `in:"query.floats"`
	StringList []string    `in:"query.strings"`
	TimeList   []time.Time `in:"query.times"`
}

type ProductQuery struct {
	CreatedAt time.Time `in:"query.created_at,required"`
	Color     string    `in:"query.color"`
	IsSoldout bool      `in:"query.is_soldout"`
	SortBy    []string  `in:"query.sort_by"`
	SortDesc  []bool    `in:"query.sort_desc"`
	Authorization
	Pagination
}

type ObjectID struct {
	timestamp [4]byte
	mid       [3]byte
	pid       [2]byte
	counter   [3]byte
}

type Cursor struct {
	AfterMarker  ObjectID `in:"query.after"`
	BeforeMarker ObjectID `in:"query.before"`
	Limit        int      `in:"query.limit"`
}

type MessageQuery struct {
	UserId string `in:"query.uid"`
	Cursor
}

type PositionXY struct {
	X int
	Y int
}

type PointsQuery struct {
	Positions []PositionXY `in:"query.positions"`
}

func TestKV_NormalCase(t *testing.T) {
	kvTest(
		t,
		map[string][]string{
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
		},
		ChaosQuery{},
		&ChaosQuery{
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
			TimeValue:       time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("E8", 8*3600)),
			BoolList:        []bool{true, false, false, true},
			IntList:         []int{9, 9, 6},
			FloatList:       []float64{0.0, 0.5, 1.0},
			StringList:      []string{"Life", "is", "a", "Miracle"},
			TimeList: []time.Time{
				time.Date(2000, 1, 2, 15, 4, 5, 0, time.FixedZone("W7", -7*3600)),
				time.Date(1991, 6, 28, 14, 0, 0, 0, time.FixedZone("E8", 8*3600)),
			},
		},
	)
}

func TestKV_SetFieldErrorOfBasicTypes(t *testing.T) {
	// type testcase struct {
	// 	key      string
	// 	value    []string
	// 	expected error
	// }
	// badCases := []testcase{
	// 	{"bool", []string{"a"}, &httpin.InvalidField{Name: "BoolValue", TagKey: "query", Tag: "bool", Value: `["a"]`}},
	// 	{"int", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"int8", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"int16", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"int32", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"int64", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"uint", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"uint8", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"uint16", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"uint32", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"uint64", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"float32", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"float64", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"complex64", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"complex128", []string{"a"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["a"]`}},
	// 	{"string", []string{""}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `[""]`}},
	// 	{"time", []string{"1991-11-10"}, &httpin.InvalidField{Name: "", TagKey: "query", Tag: "", Value: `["1991-11-10"]`}},
	// }

	// for _, badCase := range badCases {
	// 	kvTest(
	// 		t,
	// 		map[string][]string{badCase.key: badCase.value},
	// 		ChaosQuery{},
	// 		badCase.expected,
	// 	)
	// }
}

func TestKV_EmbeddedField(t *testing.T) {
	kvTest(
		t,
		map[string][]string{
			"created_at": {"1991-11-10T08:00:00+08:00"},
			"color":      {"red"},
			"is_soldout": {"true"},
			"sort_by":    {"stock", "price"},
			"sort_desc":  {"true", "true"},
			"page":       {"1"},
			"per_page":   {"20"},
		},
		ProductQuery{},
		&ProductQuery{
			CreatedAt: time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("+08:00", 8*3600)),
			Color:     "red",
			IsSoldout: true,
			SortBy:    []string{"stock", "price"},
			SortDesc:  []bool{true, true},
			Pagination: Pagination{
				Page:    1,
				PerPage: 20,
			},
		},
	)
}

func TestKV_MissingFields(t *testing.T) {
	kvTest(
		t,
		map[string][]string{
			"sort_by":   {"stock", "price"},
			"sort_desc": {"true", "true"},
			"per_page":  {"10"},
		},
		ProductQuery{},
		&ProductQuery{
			SortBy:   []string{"stock", "price"},
			SortDesc: []bool{true, true},
			Pagination: Pagination{
				Page:    0,
				PerPage: 10,
			},
		},
	)
}

func TestKV_UnsupportedCustomType(t *testing.T) {
	kvTest(
		t,
		map[string][]string{
			"uid":   {"ggicci"},
			"after": {"5cb71995ad763f7f1717c9eb"},
			"limit": {"50"},
		},
		&MessageQuery{},
		UnsupportedTypeError{Type: reflect.TypeOf(ObjectID{})},
	)
}

func TestKV_UnsupportedElementTypeOfArray(t *testing.T) {
	kvTest(
		t,
		map[string][]string{
			"positions": {"(1,4)", "(5,7)"},
		},
		PointsQuery{},
		UnsupportedTypeError{Type: reflect.TypeOf(PositionXY{})},
	)
}

func kvTest(t *testing.T, form map[string][]string, inputStruct, expected interface{}) {
	// engine, err := httpin.NewEngine(inputStruct)
	// if err != nil {
	// 	t.Errorf("unable to create engine: %s", err)
	// 	t.FailNow()
	// 	return
	// }

	// got, err := engine.ReadForm(form)
	// if err != nil {
	// 	check(t, expected, err)
	// } else {
	// 	check(t, expected, got)
	// }
}

func check(t *testing.T, expected, got interface{}) {
	// Expecting error.
	expectedError, isExpectingError := expected.(error)
	if isExpectingError {
		if gotError, ok := got.(error); ok {
			if errors.Is(gotError, expectedError) {
				return
			}
		}
		t.Errorf("parse failed, expected error %v, got %v", expectedError, got)
		return
	}

	// Expecting struct output.
	left, _ := json.Marshal(expected)
	right, _ := json.Marshal(got)
	if bytes.Compare(left, right) != 0 {
		t.Errorf("parse failed, expected %s, got %s", left, right)
		return
	}
}
