package httpin_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ggicci/httpin"
)

func TestResolver_Build(t *testing.T) {
	resolver, err := httpin.BuildFieldResolver(reflect.TypeOf(ProductQuery{}))
	if err != nil {
		t.Error("build resolver failed:", err)
		t.Fail()
	}
	jsonContent, _ := json.Marshal(resolver)
	t.Log(string(jsonContent))
}
