package httpin

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldResolver(t *testing.T) {
	Convey("Build resolver tree", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(ProductQuery{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)
		r, _ := http.NewRequest("GET", "https://example.com", nil)
		r.Form = make(url.Values)
		r.Form.Set("created_at", time.Now().Format(time.RFC3339))
		r.Form.Set("color", "red")
		r.Form.Set("is_soldout", "true")
		r.Form.Add("sort_by", "id")
		r.Form.Add("sort_by", "quantity")
		r.Form.Add("sort_desc", "0")
		r.Form.Add("sort_desc", "true")
		r.Form.Set("page", "1")
		r.Form.Set("per_page", "20")
		r.Header.Set("x-api-token", "cad979df-5e40-4bfd-b31d-f870ca2c14ea")
		rv, err := resolver.resolve(r)
		So(err, ShouldBeNil)
		So(rv.Elem().Interface(), ShouldHaveSameTypeAs, ProductQuery{})
		bs, _ := json.Marshal(rv.Interface())
		t.Logf("ProductQuery: %s\n", bs)
	})
}

func TestResolverWithMissingRequiredField(t *testing.T) {
	Convey("A resolver with a missing required field", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(ProductQuery{}))
		r, _ := http.NewRequest("GET", "https://example.com", nil)
		_, err = resolver.resolve(r)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual,"invalid field \"created_at\": missing required field")
	})
}
