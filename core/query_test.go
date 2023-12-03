package core

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/ggicci/httpin/internal"
	"github.com/stretchr/testify/assert"
)

func TestDirectiveQuery_Decode(t *testing.T) {
	type SearchQuery struct {
		Query      string `in:"query=q;required"`
		PageNumber int    `in:"query=p"`
		PageSize   int    `in:"query=page_size"`
	}

	r, _ := http.NewRequest("GET", "/?q=doggy&p=2&page_size=5", nil)
	expected := &SearchQuery{
		Query:      "doggy",
		PageNumber: 2,
		PageSize:   5,
	}

	co, err := New(SearchQuery{})
	assert.NoError(t, err)
	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.(*SearchQuery))
}

func TestDirectiveQuery_Encode(t *testing.T) {
	type SearchQuery struct {
		Name    string  `in:"query=name"`
		Age     int     `in:"query=age"`
		Enabled bool    `in:"query=enabled"`
		Price   float64 `in:"query=price"`

		NameList []string `in:"query=name_list[]"`
		AgeList  []int    `in:"query=age_list[]"`

		NamePointer *string `in:"query=name_pointer"`
		AgePointer  *int    `in:"query=age_pointer"`
	}
	query := &SearchQuery{
		Name:     "cupcake",
		Age:      12,
		Enabled:  true,
		Price:    6.28,
		NameList: []string{"apple", "banana", "cherry"},
		AgeList:  []int{1, 2, 3},
		NamePointer: func() *string {
			s := "pointer cupcake"
			return &s
		}(),
		AgePointer: func() *int {
			i := 19
			return &i
		}(),
	}

	co, err := New(SearchQuery{})
	assert.NoError(t, err)
	req, err := co.NewRequest("GET", "/pets", query)
	assert.NoError(t, err)

	expected, _ := http.NewRequest("GET", "/pets", nil)
	expectedQuery := make(url.Values)
	expectedQuery.Set("name", query.Name)                 // query.Name
	expectedQuery.Set("age", "12")                        // query.Age
	expectedQuery.Set("enabled", "true")                  // query.Enabled
	expectedQuery.Set("price", "6.28")                    // query.Price
	expectedQuery["name_list[]"] = query.NameList         // query.NameList
	expectedQuery["age_list[]"] = []string{"1", "2", "3"} // query.AgeList
	expectedQuery.Set("name_pointer", *query.NamePointer) // query.NamePointer
	expectedQuery.Set("age_pointer", "19")                // query.PointerAge
	expected.URL.RawQuery = expectedQuery.Encode()
	assert.Equal(t, expected, req)
}

type Location struct {
	Latitude  float64
	Longitude float64
}

func (l Location) ToString() (string, error) {
	return fmt.Sprintf("%f,%f", l.Latitude, l.Longitude), nil
}

type LocationTextMarshalerImpl Location

func (l LocationTextMarshalerImpl) MarshalText() ([]byte, error) {
	if s, err := (Location)(l).ToString(); err != nil {
		return nil, err
	} else {
		return []byte("MarshalText:" + s), nil
	}
}

func TestDirectiveQuery_Encode_useMarshalerInterfaces(t *testing.T) {
	type SearchQuery struct {
		L0     *Location                  `in:"query=l0"`
		L2     *LocationTextMarshalerImpl `in:"query=l2"`
		Radius int                        `in:"query=radius"`
	}

	query := &SearchQuery{
		L0: &Location{
			Latitude:  1.234,
			Longitude: 5.678,
		},
		L2: &LocationTextMarshalerImpl{
			Latitude:  1.234,
			Longitude: 5.678,
		},
		Radius: 1000,
	}

	co, err := New(SearchQuery{})
	assert.NoError(t, err)
	req, err := co.NewRequest("GET", "/pets", query)
	assert.NoError(t, err)

	expected, _ := http.NewRequest("GET", "/pets", nil)
	expectedQuery := make(url.Values)
	expectedQuery.Set("l0", "1.234000,5.678000")
	expectedQuery.Set("l2", "MarshalText:1.234000,5.678000")
	expectedQuery.Set("radius", "1000")
	expected.URL.RawQuery = expectedQuery.Encode()
	assert.Equal(t, expected, req)
}

func TestDirectiveQuery_Encode_ErrUnsupportedType(t *testing.T) {
	type SearchQuery struct {
		Map map[string]string `in:"query=map"` // unsupported type: map
	}

	co, err := New(SearchQuery{})
	assert.NoError(t, err)
	_, err = co.NewRequest("GET", "/pets", &SearchQuery{})
	assert.ErrorIs(t, err, internal.ErrUnsupportedType)
}
