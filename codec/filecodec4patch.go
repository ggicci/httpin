package codec

import (
	"io"
	"reflect"
)

// FileCodec4PatchField makes patch.Field[T] implement FileCodec as long as
// T implements FileCodec. It is used to eliminate the effort of implementing
// FileCodec for patch.Field[T] for every type T.
type FileCodec4PatchField struct {
	Value reflect.Value // of patch.Field[T]
	codec FileCodec
}

func NewFileCodec4PatchField(rv reflect.Value) (*FileCodec4PatchField, error) {
	fcodec, err := NewFileCodec(rv.FieldByName("Value"))
	if err != nil {
		return nil, err
	} else {
		return &FileCodec4PatchField{
			Value: rv,
			codec: fcodec,
		}, nil
	}
}

func (fc *FileCodec4PatchField) Filename() string {
	return fc.codec.Filename()
}

func (fc *FileCodec4PatchField) MarshalFile() (io.ReadCloser, error) {
	return fc.codec.MarshalFile()
}

func (fc *FileCodec4PatchField) UnmarshalFile(fh FileObject) error {
	if err := fc.codec.UnmarshalFile(fh); err != nil {
		return err
	} else {
		fc.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}
