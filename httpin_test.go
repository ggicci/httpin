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

func TestForm_NormalCase(t *testing.T) {
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
