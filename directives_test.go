package httpin

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type NoopDirective struct{}

func (d *NoopDirective) Execute(ctx *DirectiveContext) error {
	return nil
}

func TestDirectives(t *testing.T) {
	Convey("Register duplicate executor", t, func() {
		So(func() { RegisterDirectiveExecutor("noop", &NoopDirective{}, nil) }, ShouldNotPanic)
		So(func() { RegisterDirectiveExecutor("noop", &NoopDirective{}, nil) }, ShouldPanic)
	})

	Convey("Register nil executor", t, func() {
		So(func() { RegisterDirectiveExecutor("whatever", nil, nil) }, ShouldPanic)
	})

	Convey("Replace an executor", t, func() {
		So(func() { ReplaceDirectiveExecutor("noop", &NoopDirective{}, nil) }, ShouldNotPanic)
		So(func() { ReplaceDirectiveExecutor("noop", &NoopDirective{}, nil) }, ShouldNotPanic)
	})
}
