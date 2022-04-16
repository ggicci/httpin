package httpin

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type ThingWithDefaultValues struct {
	Page      int      `in:"form=page;default=1"`
	PerPage   int      `in:"form=per_page;default=20"`
	StateList []string `in:"form=state;default=pending,in_progress,failed"`
}

func TestDirectiveDefault(t *testing.T) {
	Convey("Default directive should set unsetted fields", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Form = url.Values{
			"page":  {"7"},
			"state": {},
		}
		expected := &ThingWithDefaultValues{
			Page:      7,
			PerPage:   20,
			StateList: []string{"pending", "in_progress", "failed"},
		}
		core, err := New(ThingWithDefaultValues{})
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, expected)
	})
}
