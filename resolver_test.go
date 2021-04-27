package httpin_test

import (
	"reflect"
	"testing"

	"github.com/ggicci/httpin"
)

func TestResolver_Build(t *testing.T) {
	resolver, err := httpin.BuildTypeResolver(reflect.TypeOf(ProductQuery{}))
	if err != nil {
		t.Error("build type resolver failed")
		t.Fail()
	}
	debug := resolver.DumpTree()
	if len(debug) == 0 {
		t.Error("empty tree")
		t.Fail()
	}
	t.Log(debug)
}
