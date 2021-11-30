package httpin

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {
	Convey("Override default error handler", t, func() {
		engine, _ := New(ProductQuery{})
		So(engine.getErrorHandler(), ShouldEqual, globalCustomErrorHandler)

		engine, _ = New(ProductQuery{}, WithErrorHandler(CustomErrorHandler))
		So(engine.getErrorHandler(), ShouldEqual, CustomErrorHandler)
	})

	Convey("Can't create engine with nil custom error handler", t, func() {
		_, err := New(ProductQuery{}, WithErrorHandler(nil))
		So(errors.Is(err, ErrNilErrorHandler), ShouldBeTrue)
	})
}
