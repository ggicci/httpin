package codec

import (
	"errors"
	"fmt"
	"reflect"
)

type FileSliceCodec4Slice struct {
	Value reflect.Value
}

func NewFileSliceCodec4Slice(rv reflect.Value) (*FileSliceCodec4Slice, error) {
	if !rv.CanAddr() {
		return nil, errors.New("unaddressable value")
	}
	return &FileSliceCodec4Slice{Value: rv}, nil
}

func (fsc *FileSliceCodec4Slice) ToFileSlice() ([]FileMarshaler, error) {
	var files = make([]FileMarshaler, fsc.Value.Len())
	for i := 0; i < fsc.Value.Len(); i++ {
		if fcodec, err := NewFileCodec(fsc.Value.Index(i)); err != nil {
			return nil, fmt.Errorf("cannot create FileCodec at index %d: %w", i, err)
		} else {
			files[i] = fcodec
		}
	}
	return files, nil
}

func (fsc *FileSliceCodec4Slice) FromFileSlice(fhs []FileObject) error {
	fsc.Value.Set(reflect.MakeSlice(fsc.Value.Type(), len(fhs), len(fhs)))
	for i, fh := range fhs {
		fcodec, err := NewFileCodec(fsc.Value.Index(i))
		if err != nil {
			return fmt.Errorf("cannot create FileCodec at index %d: %w", i, err)
		}
		if err := fcodec.UnmarshalFile(fh); err != nil {
			return fmt.Errorf("cannot unmarshal file %q at index %d: %w", fh.Filename(), i, err)
		}
	}
	return nil
}
