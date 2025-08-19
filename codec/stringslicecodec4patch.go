package codec

import "reflect"

// StringSliceCodec4PatchField makes patch.Field[T] implement StringSliceCodec
// as long as T implements StringCodec. It is used to eliminate the effort of
// implementing StringSliceCodec for patch.Field[T] for every type T.
type StringSliceCodec4PatchField struct {
	Value reflect.Value // of patch.Field[T]
	codec StringSliceCodec
}

func (ns *Namespace) NewStringSliceCodec4PatchField(rv reflect.Value, adaptor StringCodecAdaptor) (*StringSliceCodec4PatchField, error) {
	codec, err := ns.NewStringSliceCodec(rv.FieldByName("Value"), adaptor)
	if err != nil {
		return nil, err
	} else {
		return &StringSliceCodec4PatchField{
			Value: rv,
			codec: codec,
		}, nil
	}
}

func (ssc *StringSliceCodec4PatchField) ToStringSlice() ([]string, error) {
	if ssc.Value.FieldByName("Valid").Bool() {
		return ssc.codec.ToStringSlice()
	} else {
		return []string{}, nil
	}
}

func (ssc *StringSliceCodec4PatchField) FromStringSlice(values []string) error {
	if err := ssc.codec.FromStringSlice(values); err != nil {
		return err
	} else {
		ssc.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}
