package core

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/ggicci/owl"
)

type contextKey int

const (
	// CtxRequest is the key to get the HTTP request value (of *http.Request)
	// from DirectiveRuntime.Context. The HTTP request value is injected by
	// httpin to the context of DirectiveRuntime before executing the directive.
	// See Core.Decode() for more details.
	//
	// Use DirectiveRuntime.GetRequest() to get the HTTP request.
	CtxRequest contextKey = iota

	// Use DirectiveRuntime.GetRequestBuilder() to get the request builder.
	CtxRequestBuilder

	// CtxCustomCodec is the context key used to retrieve a custom codec for a field
	// from Resolver.Context. This codec is specified by the "codec" directive.
	//
	// During the resolver-building phase, the "codec" directive is stripped from the field,
	// and the corresponding codec (looked up by name) is stored in Resolver.Context using this key.
	//
	// Example:
	//
	//    type GreetInput struct {
	//        Message string `in:"codec=custom"`
	//    }
	//
	// In this example, the codec named "custom" will be associated with the resolver for the
	// Message field via this context key.
	//
	// Use DirectiveRuntime.GetCustomCodec() to get the codec.
	CtxCustomCodec

	// CtxFieldSet is used by executors to tell whether a field has been set. When
	// multiple executors were applied to a field, if the field value were set
	// by a former executor, the latter executors MAY skip running by consulting
	// this context value.
	//
	// Use DirectiveRuntime.IsFiledSet() to get the value.
	//
	// Use DirectiveRuntime.MarkFieldSet(true/false) to set the value.
	CtxFieldSet

	// CtxNamespace indicates the httpin namespace to which the resources belong.
	//
	// Use DirectiveRuntime.GetNamespace() to get the httpin namespace.
	CtxNamespace
)

// DirectiveRuntime is the runtime of a directive execution. It wraps owl.DirectiveRuntime,
// providing some additional helper methods particular to httpin.
//
// See owl.DirectiveRuntime for more details.
type DirectiveRuntime owl.DirectiveRuntime

// Get the httpin namespace.
func (rtm *DirectiveRuntime) GetNamespace() *Namespace {
	if ns := rtm.Context.Value(CtxNamespace); ns != nil {
		return ns.(*Namespace)
	}
	panic("namespace must be set")
}

// GetRequest returns the HTTP request value from the context of
// DirectiveRuntime. This is useful for executors that need to access the HTTP
// request, such as "query", "header", "cookie", etc.
func (rtm *DirectiveRuntime) GetRequest() *http.Request {
	if req := rtm.Context.Value(CtxRequest); req != nil {
		return req.(*http.Request)
	}
	return nil
}

// GetRequestBuilder returns the RequestBuilder from the context of
// DirectiveRuntime. The RequestBuilder is used to build the HTTP request
// from the directive arguments. It is useful for executors that need to
// build the HTTP request, such as "query", "header", "cookie", etc.
func (rtm *DirectiveRuntime) GetRequestBuilder() *RequestBuilder {
	if rb := rtm.Context.Value(CtxRequestBuilder); rb != nil {
		return rb.(*RequestBuilder)
	}
	return nil
}

// GetCustomCodec returns the custom codec bound to the field by the "codec", "coder",
// "decoder" directives.
func (rtm *DirectiveRuntime) GetCustomCodec() *NamedStringCodecAdaptor {
	if info := rtm.Resolver.Context.Value(CtxCustomCodec); info != nil {
		return info.(*NamedStringCodecAdaptor)
	} else {
		return nil
	}
}

// Deprecated: Use GetCustomCodec instead.
func (rtm *DirectiveRuntime) GetCustomCoder() *NamedStringCodecAdaptor {
	return rtm.GetCustomCodec()
}

// IsFieldSet checks whether the field has been set by a former executor.
// If the field has been set, the latter executors MAY skip running.
// This is useful when multiple executors are applied to a field, and you want
// to avoid running the latter executors if the field has been set by a former
// executor. For example:
//
//	token string `in:"query=access_token;header=x-api-token"
//
// If the "query" executor has set the field (i.e., got a value of access_token
// key from the querystring), the "header" executor can skip running, because
// the field is already set.
func (rtm *DirectiveRuntime) IsFieldSet() bool {
	return rtm.Context.Value(CtxFieldSet) == true
}

// MarkFieldSet marks the field as set. This is used by executors to tell
// that the field has been set. The latter executors can consult this value
// to decide whether to skip running.
func (rtm *DirectiveRuntime) MarkFieldSet(value bool) {
	rtm.Context = context.WithValue(rtm.Context, CtxFieldSet, value)
}

// SetValue sets the value of the field that being wrapped by this
// DirectiveRuntime. It is useful for users who is implementing custom
// directives and need to set the value of a field directly. This helper method
// will check the type of the value and ensure it is assignable to the field's
// type and throws an appropriate error on failure.
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
			ErrFieldTypeMismatch, reflect.TypeOf(value), targetType)
	}

	rtm.Value.Elem().Set(newValue)
	return nil
}
