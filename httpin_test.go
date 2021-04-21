package httpin_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/ggicci/httpin"
)

type TestCase struct {
	Title         string
	InputType     interface{}
	InputForm     url.Values
	InputBody     io.Reader
	Expected      interface{}
	ExpectedError error
}

func (c *TestCase) Check(t *testing.T, got interface{}, err error) {
	if !errors.Is(err, c.ExpectedError) {
		t.Errorf("parse failed, expect error %v, got %v", c.ExpectedError, err)
		return
	}

	// Compare got vs. expected.
	left, _ := json.Marshal(c.Expected)
	right, _ := json.Marshal(got)
	if bytes.Compare(left, right) != 0 {
		t.Errorf("parse failed, expected %s, got %s", left, right)
		return
	}
}

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

func TestExtractingQueryParameters(t *testing.T) {
	testcases := []*TestCase{
		{
			Title:     "normal case",
			InputType: ProductQuery{},
			InputForm: url.Values{
				"created_at": {"1991-11-10T08:00:00+08:00"},
				"color":      {"red"},
				"is_soldout": {"true"},
				"sort_by":    {"stock", "price"},
				"sort_desc":  {"true", "true"},
				"page":       {"1"},
				"per_page":   {"20"},
			},
			Expected: &ProductQuery{
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
		},
		{
			Title:     "unsupported custom type",
			InputType: &MessageQuery{},
			InputForm: url.Values{
				"uid":   {"ggicci"},
				"after": {"5cb71995ad763f7f1717c9eb"},
				"limit": {"50"},
			},
			ExpectedError: httpin.UnsupportedType("ObjectID"),
		},
	}

	for _, c := range testcases {
		engine, err := httpin.NewEngine(c.InputType)
		if err != nil {
			t.Errorf("unable to create engine: %s", err)
			t.FailNow()
		}

		got, err := engine.ReadForm(c.InputForm)
		c.Check(t, got, err)
	}
}
