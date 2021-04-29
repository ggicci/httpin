package httpin

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldResolver(t *testing.T) {
	Convey("Build resolver tree", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(ProductQuery{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)

		// TODO(ggicci): add more checks and tests
	})
}
