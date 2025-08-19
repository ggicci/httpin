package codec

import (
	"reflect"
)

type FileSliceCodec interface {
	ToFileSlice() ([]FileMarshaler, error)
	FromFileSlice([]FileObject) error
}

func NewFileSliceCodec(rv reflect.Value) (FileSliceCodec, error) {
	if IsPatchField(rv.Type()) {
		return NewFileSliceCodec4PatchField(rv)
	}

	if isSliceType(rv.Type()) {
		return NewFileSliceCodec4Slice(rv)
	} else {
		return NewFileSliceCodec4FileCodec(rv)
	}
}
