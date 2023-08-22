package httpin

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {
	Convey("Override default error handler", t, func() {
		core, _ := New(ProductQuery{})
		So(core.getErrorHandler(), ShouldEqual, globalCustomErrorHandler)

		core, _ = New(ProductQuery{}, WithErrorHandler(CustomErrorHandler))
		So(core.getErrorHandler(), ShouldEqual, CustomErrorHandler)
	})

	Convey("Can't create core with nil custom error handler", t, func() {
		_, err := New(ProductQuery{}, WithErrorHandler(nil))
		So(errors.Is(err, ErrNilErrorHandler), ShouldBeTrue)
	})

	Convey("Use invalid max memory", t, func() {
		_, err := New(ProductQuery{}, WithMaxMemory(100))
		So(errors.Is(err, ErrMaxMemoryTooSmall), ShouldBeTrue)
	})
}
