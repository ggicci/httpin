package httpin

import (
	"fmt"
	"mime/multipart"
	"reflect"
)

func init() {
	RegisterTypeDecoder(reflect.TypeOf(File{}), FileTypeDecoderFunc(DecodeFile))
}

type File struct {
	multipart.File
	Header *multipart.FileHeader
}

func DecodeFile(meta *multipart.FileHeader) (interface{}, error) {
	if meta == nil {
		return nil, ErrNilFile
	}

	file, err := meta.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	inFile := File{
		File:   file,
		Header: meta,
	}

	return inFile, nil
}
