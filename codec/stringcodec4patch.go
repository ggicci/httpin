package codec

import (
	"errors"
	"fmt"
	"reflect"
)

// StringCodec4PatchField makes patch.Field[T] implement StringCodec as long as
// T implements StringCodec. It is used to eliminate the effort of implementing
// StringCodec for patch.Field[T] for every type T.
type StringCodec4PatchField struct {
	Value reflect.Value // of patch.Field[T]
	codec StringCodec
}

func (ns *Namespace) NewStringCodec4PatchField(rv reflect.Value, adapt StringCodecAdaptor) (*StringCodec4PatchField, error) {
	StringCodec, err := ns.NewStringCodec(rv.FieldByName("Value"), adapt)
	if err != nil {
		return &StringCodec4PatchField{}, fmt.Errorf("cannot create StringCodec for PatchField: %w", err)
	} else {
		return &StringCodec4PatchField{
			Value: rv,
			codec: StringCodec,
		}, nil
	}
}

func (w *StringCodec4PatchField) ToString() (string, error) {
	if w.Value.FieldByName("Valid").Bool() {
		return w.codec.ToString()
	} else {
		return "", errors.New("invalid value") // when Valid is false
	}
}

// FromString sets the value of the wrapped patch.Field[T] from the given
// string. It returns an error if the given string is not valid. And leaves the
// original value of both Value and Valid unchanged. On the other hand, if no
// error occurs, it sets Valid to true.
func (w *StringCodec4PatchField) FromString(s string) error {
	if err := w.codec.FromString(s); err != nil {
		return err
	} else {
		w.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}
