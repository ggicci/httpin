package codec

import "reflect"

// FileSliceCodec4PatchField makes patch.Field[T] implement FileSliceCodec as long as
// T implements FileSliceCodec. It is used to eliminate the effort of implementing
// FileSliceCodec for patch.Field[T] for every type T.
type FileSliceCodec4PatchField struct {
	Value reflect.Value // of patch.Field[T]
	codec FileSliceCodec
}

func NewFileSliceCodec4PatchField(rv reflect.Value) (*FileSliceCodec4PatchField, error) {
	fileSliceCodec, err := NewFileSliceCodec(rv.FieldByName("Value"))
	if err != nil {
		return nil, err
	} else {
		return &FileSliceCodec4PatchField{
			Value: rv,
			codec: fileSliceCodec,
		}, nil
	}
}

func (fsc *FileSliceCodec4PatchField) ToFileSlice() ([]FileMarshaler, error) {
	if fsc.Value.FieldByName("Valid").Bool() {
		return fsc.codec.ToFileSlice()
	} else {
		return []FileMarshaler{}, nil
	}
}

func (fsc *FileSliceCodec4PatchField) FromFileSlice(fhs []FileObject) error {
	if err := fsc.codec.FromFileSlice(fhs); err != nil {
		return err
	} else {
		fsc.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}
