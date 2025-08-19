package codec

import "reflect"

type FileSliceCodec4FileCodec struct{ FileCodec }

func NewFileSliceCodec4FileCodec(rv reflect.Value) (*FileSliceCodec4FileCodec, error) {
	if fcodec, err := NewFileCodec(rv); err != nil {
		return nil, err
	} else {
		return &FileSliceCodec4FileCodec{fcodec}, nil
	}
}

func (fsc *FileSliceCodec4FileCodec) ToFileSlice() ([]FileMarshaler, error) {
	return []FileMarshaler{fsc.FileCodec}, nil
}

func (fsc *FileSliceCodec4FileCodec) FromFileSlice(files []FileObject) error {
	if len(files) > 0 {
		return fsc.UnmarshalFile(files[0])
	}
	return nil
}
