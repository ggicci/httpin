package httpin_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/ggicci/httpin"
)

type Pagination struct {
	Page    int `query:"page"`
	PerPage int `query:"per_page"`
}

// ChaosQuery is designed to make the normal case test coverage higher.
type ChaosQuery struct {
	// Basic Types
	BoolValue       bool       `query:"bool"`
	IntValue        int        `query:"int"`
	Int8Value       int8       `query:"int8"`
	Int16Value      int16      `query:"int16"`
	Int32Value      int32      `query:"int32"`
	Int64Value      int64      `query:"int64"`
	UintValue       uint       `query:"uint"`
	Uint8Value      uint8      `query:"uint8"`
	Uint16Value     uint16     `query:"uint16"`
	Uint32Value     uint32     `query:"uint32"`
	Uint64Value     uint64     `query:"uint64"`
	Float32Value    float32    `query:"float32"`
	Float64Value    float64    `query:"float64"`
	Complex64Value  complex64  `query:"complex64"`
	Complex128Value complex128 `query:"complex128"`
	StringValue     string     `query:"string"`

	// Time Type
	TimeValue time.Time `query:"time"`

	// Array
	BoolList   []bool      `query:"bools"`
	IntList    []int       `query:"ints"`
	FloatList  []float64   `query:"floats"`
	StringList []string    `query:"strings"`
	TimeList   []time.Time `query:"times"`
}

type ProductQuery struct {
	CreatedAt time.Time `query:"created_at"`
	Color     string    `query:"color"`
	IsSoldout bool      `query:"is_soldout"`
	SortBy    []string  `query:"sort_by"`
	SortDesc  []bool    `query:"sort_desc"`
	Pagination
}

type ObjectID struct {
	timestamp [4]byte
	mid       [3]byte
	pid       [2]byte
	counter   [3]byte
}

type Cursor struct {
	AfterMarker  ObjectID `query:"after"`
	BeforeMarker ObjectID `query:"before"`
	Limit        int      `query:"limit"`
}

type MessageQuery struct {
	UserId string `query:"uid"`
	Cursor
}

type PositionXY struct {
	X int
	Y int
}

type PointsQuery struct {
	Positions []PositionXY `query:"positions"`
}

func TestForm_NormalCase(t *testing.T) {
	parseFormAndCheck(
		t,
		url.Values{
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

func TestForm_EmbeddedField(t *testing.T) {
	parseFormAndCheck(
		t,
		url.Values{

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

func TestForm_UnsupportedCustomType(t *testing.T) {
	parseFormAndCheck(
		t,
		url.Values{
			"uid":   {"ggicci"},
			"after": {"5cb71995ad763f7f1717c9eb"},
			"limit": {"50"},
		},
		&MessageQuery{},
		httpin.UnsupportedType("ObjectID"),
	)
}

func TestForm_UnsupportedElementTypeOfArray(t *testing.T) {
	parseFormAndCheck(
		t,
		url.Values{
			"positions": {"(1,4)", "(5,7)"},
		},
		PointsQuery{},
		httpin.UnsupportedType("PositionXY"),
	)
}

func parseFormAndCheck(t *testing.T, form url.Values, inputStruct, expected interface{}) {
	engine, err := httpin.NewEngine(inputStruct)
	if err != nil {
		t.Errorf("unable to create engine: %s", err)
		t.FailNow()
		return
	}

	got, err := engine.ReadForm(form)
	if err != nil {
		check(t, expected, err)
	} else {
		check(t, expected, got)
	}
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
