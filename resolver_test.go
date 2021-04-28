package httpin_test

import (
	"reflect"
	"testing"

	"github.com/ggicci/httpin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldResolver(t *testing.T) {
	Convey("Build FieldResolver normally", t, func() {
		resolver, err := httpin.BuildFieldResolver(reflect.TypeOf(ProductQuery{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)

		// TODO(ggicci): add more checks and tests
	})
}

func TestResolver_BuildNonStructType(t *testing.T) {
	Convey("Build FieldResolver with non-struct type", t, func() {
		var Name string
		resolver, err := httpin.BuildFieldResolver(reflect.TypeOf(Name))
		So(err, ShouldBeError)
		So(resolver, ShouldBeNil)
	})
}
