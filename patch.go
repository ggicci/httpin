package httpin

import (
	"github.com/ggicci/owl"
)

type validSetter interface {
	SetValid(bool)
}

func setField(rtm *DirectiveRuntime) error {
	if rtm.Context.Value(FieldSet) == true {
		elem := rtm.Value.Elem().Interface()
		if setter, ok := elem.(validSetter); ok {
			setter.SetValid(true)
		}
	}
	return nil
}

func appendSetFieldDirective(r *owl.Resolver) error {
	r.Directives = append(r.Directives, owl.NewDirective("_setfield"))
	return nil
}
