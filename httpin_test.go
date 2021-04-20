package httpin_test

import (
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

	// TODO(ggicci): compare got vs. expected
}

type ProductQuery struct {
	CreatedAt time.Time `query:"created_at"`
	Color     string    `query:"color"`
	IsSoldout bool      `query:"is_soldout"`
	SortBy    []string  `query:"sort_by"`
	SortDesc  []bool    `query:"sort_desc"`
	Page      int       `query:"page"`
	PerPage   int       `query:"per_per"`
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
				"created_at": {"2020-01-02"},
				"color":      {"red"},
				"is_soldout": {"true"},
			},
		},
	}

	for _, c := range testcases {
		got, err := engine.ReadForm(c.InputForm)
		c.Check(t, got, err)
	}
}
