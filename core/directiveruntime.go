package core

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/owl"
)

type contextKey int

const (
	// CtxRequest is the key to get the HTTP request value (of *http.Request)
	// from DirectiveRuntime.Context. The HTTP request value is injected by
	// httpin to the context of DirectiveRuntime before executing the directive.
	// See Core.Decode() for more details.
	CtxRequest contextKey = iota

	CtxRequestBuilder

	// CtxCustomDecoder is the key to get the custom decoder for a field from
	// Resolver.Context. Which is specified by the "decoder" directive.
	// During resolver building phase, the "decoder" directive will be removed
	// from the resolver, and the targeted decoder by name will be put into
	// Resolver.Context with this key. e.g.
	//
	//    type GreetInput struct {
	//        Message string `httpin:"decoder=custom"`
	//    }
	// For the above example, the decoder named "custom" will be put into the
	// resolver of Message field with this key.
	CtxCustomDecoder

	// CtxCustomEncoder works like ctxCustomDecoder, but for encoder.
	CtxCustomEncoder

	// CtxFieldSet is used by executors to tell whether a field has been set. When
	// multiple executors were applied to a field, if the field value were set
	// by a former executor, the latter executors MAY skip running by consulting
	// this context value.
	CtxFieldSet
)

// DirectiveRuntime is the runtime of a directive execution. It wraps owl.DirectiveRuntime,
// providing some additional helper methods particular to httpin.
//
// See owl.DirectiveRuntime for more details.
type DirectiveRuntime owl.DirectiveRuntime

func (rtm *DirectiveRuntime) GetRequest() *http.Request {
	if req := rtm.Context.Value(CtxRequest); req != nil {
		return req.(*http.Request)
	}
	return nil
}

func (rtm *DirectiveRuntime) GetRequestBuilder() *RequestBuilder {
	if rb := rtm.Context.Value(CtxRequestBuilder); rb != nil {
		return rb.(*RequestBuilder)
	}
	return nil
}

func (rtm *DirectiveRuntime) GetCustomDecoder() *NamedAnyStringableAdaptor {
	if info := rtm.Resolver.Context.Value(CtxCustomDecoder); info != nil {
		return info.(*NamedAnyStringableAdaptor)
	} else {
		return nil
	}
}

func (rtm *DirectiveRuntime) GetCustomEncoder() *NamedAnyStringableAdaptor {
	if info := rtm.Resolver.Context.Value(CtxCustomEncoder); info != nil {
		return info.(*NamedAnyStringableAdaptor)
	} else {
		return nil
	}
}

func (rtm *DirectiveRuntime) IsFieldSet() bool {
	return rtm.Context.Value(CtxFieldSet) == true
}

func (rtm *DirectiveRuntime) MarkFieldSet(value bool) {
	rtm.Context = context.WithValue(rtm.Context, CtxFieldSet, value)
}

func (rtm *DirectiveRuntime) SetValue(value any) error {
	if value == nil {
		// NOTE: should we wipe the value here? i.e. set the value to nil if necessary.
		// No case found yet, at least for now.
		return nil
	}
	newValue := reflect.ValueOf(value)
	targetType := rtm.Value.Type().Elem()

	if !newValue.Type().AssignableTo(targetType) {
		return fmt.Errorf("%w: value of type %q is not assignable to type %q",
			internal.ErrTypeMismatch, reflect.TypeOf(value), targetType)
	}

	rtm.Value.Elem().Set(newValue)
	return nil
}
