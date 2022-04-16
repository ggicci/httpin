package httpin

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDirectiveRequired(t *testing.T) {
	Convey("Required field is missing", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"color":      {"red"},
			"is_soldout": {"true"},
			"sort_by":    {"id", "quantity"},
			"sort_desc":  {"0", "true"},
			"page":       {"1"},
			"per_page":   {"20"},
		}
		core, err := New(&ProductQuery{}) // struct pointer also works
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(got, ShouldBeNil)
		So(errors.Is(err, ErrMissingField), ShouldBeTrue)

		var invalidField *InvalidFieldError
		So(errors.As(err, &invalidField), ShouldBeTrue)
		So(invalidField.Source, ShouldEqual, "required")
		So(invalidField.Value, ShouldBeNil)
	})

	Convey("Non-required fields can be absent", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"created_at": {"1991-11-10T08:00:00+08:00"},
			"is_soldout": {"true"},
			"page":       {"1"},
			"per_page":   {"20"},
		}
		expected := &ProductQuery{
			CreatedAt: time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC),
			Color:     "",
			IsSoldout: true,
			Pagination: Pagination{
				Page:    1,
				PerPage: 20,
			},
		}
		core, err := New(ProductQuery{})
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, expected)
	})
}
