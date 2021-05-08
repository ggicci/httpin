package httpin

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {
	Convey("Apply error status code", t, func() {
		engine, _ := New(ProductQuery{})
		So(engine.errorStatusCode, ShouldEqual, 422)

		engine, _ = New(ProductQuery{}, WithErrorStatusCode(400))
		So(engine.errorStatusCode, ShouldEqual, 400)
	})
}
