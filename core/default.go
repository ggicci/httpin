// directive: "default"
// https://ggicci.github.io/httpin/directives/default

package core

import (
	"mime/multipart"

	"github.com/ggicci/httpin/internal"
)

type DirectiveDefault struct{}

func (*DirectiveDefault) Decode(rtm *DirectiveRuntime) error {
	if rtm.IsFieldSet() {
		return nil // noop, the field was set by a former executor
	}

	// Transform:
	// 1. ctx.Argv -> input values
	// 2. ["default"] -> keys
	extractor := &FormExtractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: map[string][]string{
				"default": rtm.Directive.Argv,
			},
		},
	}
	return extractor.Extract("default")
}

func (*DirectiveDefault) Encode(rtm *DirectiveRuntime) error {
	if !rtm.Value.IsZero() {
		return nil // skip if the field is not empty
	}
	var adapt StringCodecAdaptor
	coder := rtm.GetCustomCodec()
	if coder != nil {
		adapt = coder.Adaptor
	}
	if stringSlicable, err := internal.NewStringSliceCodec(rtm.Value, adapt); err != nil {
		return err
	} else {
		return stringSlicable.FromStringSlice(rtm.Directive.Argv)
	}
}
