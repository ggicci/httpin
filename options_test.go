package httpin

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {
	Convey("Override default error handler", t, func() {
		engine, _ := New(ProductQuery{})
		So(engine.errorHandler, ShouldEqual, defaultErrorHandler)

		engine, _ = New(ProductQuery{}, WithErrorHandler(CustomErrorHandler))
		So(engine.errorHandler, ShouldEqual, CustomErrorHandler)
	})

	Convey("Nil handler should error", t, func() {
		_, err := New(ProductQuery{}, WithErrorHandler(nil))
		So(errors.Is(err, ErrNilErrorHandler), ShouldBeTrue)
	})
}
