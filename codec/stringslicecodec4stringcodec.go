package codec

import "reflect"

// StringSliceCodec4SingleStringCodec turns a single StringCodec into a StringSliceCodec.
type StringSliceCodec4SingleStringCodec struct{ StringCodec }

func (ns *Namespace) NewStringSliceCodec4SingleStringCodec(
	rv reflect.Value, adaptor StringCodecAdaptor) (*StringSliceCodec4SingleStringCodec, error) {
	if codec, err := ns.NewStringCodec(rv, adaptor); err != nil {
		return nil, err
	} else {
		return &StringSliceCodec4SingleStringCodec{codec}, nil
	}
}

func (ssc *StringSliceCodec4SingleStringCodec) ToStringSlice() ([]string, error) {
	if value, err := ssc.ToString(); err != nil {
		return nil, err
	} else {
		return []string{value}, nil
	}
}

func (ssc *StringSliceCodec4SingleStringCodec) FromStringSlice(values []string) error {
	if len(values) > 0 {
		return ssc.FromString(values[0])
	}
	return nil
}
