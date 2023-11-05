package core

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

// FileEncoder is the interface implemented by types that can represent a file upload.
type FileEncoder interface {
	// Encode returns the filename of the file and a io.ReadCloser for the file content.
	Encode() (string, io.ReadCloser, error)
}

func toAnyFileDecoder[T FileEncoder](fd FileDecoder[T]) FileDecoder[any] {
	if fd == nil {
		return nil
	}
	return FileDecoderFunc[any](func(fh *multipart.FileHeader) (any, error) {
		return fd.Decode(fh)
	})
}

func toFileEncoders(value reflect.Value, kind TypeKind) ([]FileEncoder, error) {
	if internal.IsNil(value) {
		return nil, nil // skip no file upload: value is nil
	}

	switch kind {
	case TypeKindT:
		return toFileEncodersOne(value)
	case TypeKindPatchT:
		if !value.FieldByName("Valid").Bool() {
			return nil, nil // skip no file upload: patch.Field.Valid is false
		}
		return toFileEncodersOne(value.FieldByName("Value"))
	case TypeKindTSlice:
		return toFileEncodersMulti(value)
	case TypeKindPatchTSlice:
		if !value.FieldByName("Valid").Bool() {
			return nil, nil // skip no file upload: patch.Field.Valid is false
		}
		return toFileEncodersMulti(value.FieldByName("Value"))
	}
	return nil, nil
}

func toFileEncodersOne(one reflect.Value) ([]FileEncoder, error) {
	if err := validateFileEncoderValue(one); err != nil {
		return nil, err
	}
	return []FileEncoder{one.Interface().(FileEncoder)}, nil
}

func toFileEncodersMulti(multi reflect.Value) ([]FileEncoder, error) {
	files := make([]FileEncoder, multi.Len())
	for i := 0; i < multi.Len(); i++ {
		if err := validateFileEncoderValue(multi.Index(i)); err != nil {
			return nil, fmt.Errorf("at index %d: %v", i, err)
		} else {
			files[i] = multi.Index(i).Interface().(FileEncoder)
		}
	}
	return files, nil
}

func validateFileEncoderValue(value reflect.Value) error {
	if internal.IsNil(value) {
		return errors.New("nil file encoder")
	}
	return nil
}
