package httpin

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectiveQuery(t *testing.T) {
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

	core, err := New(SearchQuery{})
	assert.NoError(t, err)
	got, err := core.Decode(r)
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

	core, err := New(SearchQuery{})
	assert.NoError(t, err)
	req, err := core.Encode("GET", "/pets", query)
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

	assertRequest(t, expected, req)
}

func TestDirectiveQuery_Encode_useMarshalerInterfaces(t *testing.T) {
	type SearchQuery struct {
		L0     *Location                       `in:"query=l0"`
		L1     *LocationFormValueMarshalerImpl `in:"query=l1"`
		L2     *LocationTextMarshalerImpl      `in:"query=l2"`
		Radius int                             `in:"query=radius"`
	}

	query := &SearchQuery{
		L0: &Location{
			Latitude:  1.234,
			Longitude: 5.678,
		},
		L1: &LocationFormValueMarshalerImpl{
			Latitude:  1.234,
			Longitude: 5.678,
		},
		L2: &LocationTextMarshalerImpl{
			Latitude:  1.234,
			Longitude: 5.678,
		},
		Radius: 1000,
	}

	core, err := New(SearchQuery{})
	assert.NoError(t, err)
	req, err := core.Encode("GET", "/pets", query)
	assert.NoError(t, err)

	expected, _ := http.NewRequest("GET", "/pets", nil)
	expectedQuery := make(url.Values)
	expectedQuery.Set("l0", "1.234000,5.678000")
	expectedQuery.Set("l1", "HttpinFormValue:1.234000,5.678000")
	expectedQuery.Set("l2", "MarshalText:1.234000,5.678000")
	expectedQuery.Set("radius", "1000")
	expected.URL.RawQuery = expectedQuery.Encode()

	assertRequest(t, req, expected)
}

func TestDirectiveQuery_Encode_ErrUnsupportedType(t *testing.T) {
	type SearchQuery struct {
		Map map[string]string `in:"query=map"` // unsupported type: map
	}

	core, err := New(SearchQuery{})
	assert.NoError(t, err)
	_, err = core.Encode("GET", "/pets", &SearchQuery{})
	assert.ErrorIs(t, err, errUnsupportedType)
}
