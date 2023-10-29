package httpin

import (
	"context"
	"net/http"
	"reflect"

	"github.com/ggicci/owl"
)

type contextKey int

const (
	// Input is the key to get the input object from Request.Context() injected by httpin. e.g.
	//
	//     input := r.Context().Value(httpin.Input).(*InputStruct)
	Input contextKey = iota

	// ctxRequest is the key to get the HTTP request value (of *http.Request)
	// from DirectiveRuntime.Context. The HTTP request value is injected by
	// httpin to the context of DirectiveRuntime before executing the directive.
	// See Core.Decode() for more details.
	ctxRequest

	ctxRequestBuilder

	// ctxCustomDecoder is the key to get the custom decoder for a field from
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
	ctxCustomDecoder

	// ctxCustomEncoder works like ctxCustomDecoder, but for encoder.
	ctxCustomEncoder

	// ctxFieldSet is used by executors to tell whether a field has been set. When
	// multiple executors were applied to a field, if the field value were set
	// by a former executor, the latter executors MAY skip running by consulting
	// this context value.
	ctxFieldSet
)

// DirectiveRuntime is the runtime of a directive execution. It wraps owl.DirectiveRuntime,
// providing some additional helper methods particular to httpin.
//
// See owl.DirectiveRuntime for more details.
type DirectiveRuntime owl.DirectiveRuntime

func (rtm *DirectiveRuntime) GetRequest() *http.Request {
	if req := rtm.Context.Value(ctxRequest); req != nil {
		return req.(*http.Request)
	}
	return nil
}

func (rtm *DirectiveRuntime) GetRequestBuilder() *RequestBuilder {
	if rb := rtm.Context.Value(ctxRequestBuilder); rb != nil {
		return rb.(*RequestBuilder)
	}
	return nil
}

func (rtm *DirectiveRuntime) GetCustomDecoder() (string, any) {
	if info := rtm.getCustomDecoder(); info != nil {
		return info.Name, info.Original
	} else {
		return "", nil
	}
}

func (rtm *DirectiveRuntime) getCustomDecoder() *namedDecoderInfo {
	if info := rtm.Resolver.Context.Value(ctxCustomDecoder); info != nil {
		return info.(*namedDecoderInfo)
	} else {
		return nil
	}
}

func (rtm *DirectiveRuntime) GetCustomEncoder() (string, Encoder) {
	if info := rtm.getCustomEncoder(); info != nil {
		return info.Name, info.Original
	} else {
		return "", nil
	}
}

func (rtm *DirectiveRuntime) getCustomEncoder() *namedEncoderInfo {
	if info := rtm.Resolver.Context.Value(ctxCustomEncoder); info != nil {
		return info.(*namedEncoderInfo)
	} else {
		return nil
	}
}

func (rtm *DirectiveRuntime) IsFieldSet() bool {
	return rtm.Context.Value(ctxFieldSet) == true
}

func (rtm *DirectiveRuntime) MarkFieldSet(value bool) {
	rtm.Context = context.WithValue(rtm.Context, ctxFieldSet, value)
}

func (rtm *DirectiveRuntime) SetValue(value any) error {
	if value == nil {
		// NOTE: should we wipe the value here? i.e. set the value to nil if necessary.
		// No case found yet, at lease for now.
		return nil
	}
	newValue := reflect.ValueOf(value)
	targetType := rtm.Value.Type().Elem()
	if newValue.Type().AssignableTo(targetType) {
		rtm.Value.Elem().Set(newValue)
		return nil
	}
	return invalidDecodeReturnType(targetType, reflect.TypeOf(value))
}
