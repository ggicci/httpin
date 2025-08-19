package codec

import (
	"errors"
	"fmt"
	"reflect"
)

// StringSliceCodec4Slice makes []T or [T] implement StringSliceCodec as long as
// T implements StringCodec. It is used to eliminate the effort of implementing
// StringSliceCodec for []T or [T] for every type T.
type StringSliceCodec4Slice struct {
	Value   reflect.Value // of []T or [T]
	Adaptor StringCodecAdaptor
	ns      *Namespace // Namespace to create StringCodec instances
}

func (ns *Namespace) NewStringSliceCodec4Slice(rv reflect.Value, adaptor StringCodecAdaptor) (*StringSliceCodec4Slice, error) {
	if !rv.CanAddr() {
		return nil, errors.New("unaddressable value")
	}
	return &StringSliceCodec4Slice{Value: rv, Adaptor: adaptor, ns: ns}, nil
}

func (ssc *StringSliceCodec4Slice) ToStringSlice() ([]string, error) {
	if codecs, err := ssc.makeCodecs(); err != nil {
		return nil, err
	} else {
		return codecs.ToStringSlice()
	}
}

func (ssc *StringSliceCodec4Slice) FromStringSlice(ss []string) error {
	ssc.Value.Set(reflect.MakeSlice(ssc.Value.Type(), len(ss), len(ss)))
	if codecs, err := ssc.makeCodecs(); err != nil {
		return err
	} else {
		return codecs.FromStringSlice(ss)
	}
}

func (ssc *StringSliceCodec4Slice) makeCodecs() (StringCodecs, error) {
	var codecs = make(StringCodecs, ssc.Value.Len())
	for i := 0; i < ssc.Value.Len(); i++ {
		if codec, err := ssc.ns.NewStringCodec(ssc.Value.Index(i), ssc.Adaptor); err != nil {
			return nil, fmt.Errorf("cannot create StringCodec from %q at index %d: %w", ssc.Value.Index(i), i, err)
		} else {
			codecs[i] = codec
		}
	}
	return codecs, nil
}
