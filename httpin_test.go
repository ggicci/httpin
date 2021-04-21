package httpin_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/ggicci/httpin"
)

type TestCase struct {
	InputForm url.Values
	InputBody io.Reader
	Expected  interface{}
}

func (c *TestCase) Check(t *testing.T, got interface{}, err error) {
	if err != nil {
		t.Errorf("parse error: %s", err)
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

type ProductQuery struct {
	CreatedAt time.Time `query:"created_at"`
	Color     string    `query:"color"`
	IsSoldout bool      `query:"is_soldout"`
	SortBy    []string  `query:"sort_by"`
	SortDesc  []bool    `query:"sort_desc"`
	Page      int       `query:"page"`
	PerPage   int       `query:"per_page"`
}

func TestExtractingQueryParameters(t *testing.T) {
	engine, err := httpin.NewEngine(ProductQuery{})
	if err != nil {
		t.Errorf("unable to create engine: %s", err)
		t.FailNow()
	}

	testcases := []*TestCase{
		{
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
				Page:      1,
				PerPage:   20,
			},
		},
	}

	for _, c := range testcases {
		got, err := engine.ReadForm(c.InputForm)
		c.Check(t, got, err)
	}
}
