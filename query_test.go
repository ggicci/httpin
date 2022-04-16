package httpin

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type SearchQuery struct {
	Query      string `in:"query=q;required"`
	PageNumber int    `in:"query=p"`
	PageSize   int    `in:"query=page_size"`
}

func TestDirectiveQuery(t *testing.T) {
	Convey("Get with QueryString params", t, func() {
		r, _ := http.NewRequest("GET", "/?q=doggy&p=2&page_size=5", nil)
		expected := &SearchQuery{
			Query:      "doggy",
			PageNumber: 2,
			PageSize:   5,
		}

		core, err := New(SearchQuery{})
		So(err, ShouldBeNil)
		got, err := core.Decode(r)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, expected)
	})
}
