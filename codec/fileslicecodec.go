package codec

import (
	"errors"
	"fmt"
	"reflect"
)

type FileSliceCodec interface {
	ToFileSlice() ([]FileMarshaler, error)
	FromFileSlice([]FileHeader) error
}

func NewFileSliceCodec(rv reflect.Value) (FileSliceCodec, error) {
	if IsPatchField(rv.Type()) {
		return NewFileSlicablePatchFieldWrapper(rv)
	}

	if isSliceType(rv.Type()) {
		return NewFileableSliceWrapper(rv)
	} else {
		return NewFileSlicableSingleFileableWrapper(rv)
	}
}

type FileSlicablePatchFieldWrapper struct {
	Value                 reflect.Value // of patch.Field[T]
	internalFileSliceable FileSliceCodec
}

func NewFileSlicablePatchFieldWrapper(rv reflect.Value) (*FileSlicablePatchFieldWrapper, error) {
	fileSlicable, err := NewFileSliceCodec(rv.FieldByName("Value"))
	if err != nil {
		return nil, err
	} else {
		return &FileSlicablePatchFieldWrapper{
			Value:                 rv,
			internalFileSliceable: fileSlicable,
		}, nil
	}
}

func (w *FileSlicablePatchFieldWrapper) ToFileSlice() ([]FileMarshaler, error) {
	if w.Value.FieldByName("Valid").Bool() {
		return w.internalFileSliceable.ToFileSlice()
	} else {
		return []FileMarshaler{}, nil
	}
}

func (w *FileSlicablePatchFieldWrapper) FromFileSlice(fhs []FileHeader) error {
	if err := w.internalFileSliceable.FromFileSlice(fhs); err != nil {
		return err
	} else {
		w.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}

type FileableSliceWrapper struct {
	Value reflect.Value
}

func NewFileableSliceWrapper(rv reflect.Value) (*FileableSliceWrapper, error) {
	if !rv.CanAddr() {
		return nil, errors.New("unaddressable value")
	}
	return &FileableSliceWrapper{Value: rv}, nil
}

func (w *FileableSliceWrapper) ToFileSlice() ([]FileMarshaler, error) {
	var files = make([]FileMarshaler, w.Value.Len())
	for i := 0; i < w.Value.Len(); i++ {
		if fileable, err := NewFileable(w.Value.Index(i)); err != nil {
			return nil, fmt.Errorf("cannot create Fileable at index %d: %w", i, err)
		} else {
			files[i] = fileable
		}
	}
	return files, nil
}

func (w *FileableSliceWrapper) FromFileSlice(fhs []FileHeader) error {
	w.Value.Set(reflect.MakeSlice(w.Value.Type(), len(fhs), len(fhs)))
	for i, fh := range fhs {
		fileable, err := NewFileable(w.Value.Index(i))
		if err != nil {
			return fmt.Errorf("cannot create Fileable at index %d: %w", i, err)
		}
		if err := fileable.UnmarshalFile(fh); err != nil {
			return fmt.Errorf("cannot unmarshal file %q at index %d: %w", fh.Filename(), i, err)
		}
	}
	return nil
}

type FileSlicableSingleFileableWrapper struct{ FileCodec }

func NewFileSlicableSingleFileableWrapper(rv reflect.Value) (*FileSlicableSingleFileableWrapper, error) {
	if fileable, err := NewFileable(rv); err != nil {
		return nil, err
	} else {
		return &FileSlicableSingleFileableWrapper{fileable}, nil
	}
}

func (w *FileSlicableSingleFileableWrapper) ToFileSlice() ([]FileMarshaler, error) {
	return []FileMarshaler{w.FileCodec}, nil
}

func (w *FileSlicableSingleFileableWrapper) FromFileSlice(files []FileHeader) error {
	if len(files) > 0 {
		return w.UnmarshalFile(files[0])
	}
	return nil
}
