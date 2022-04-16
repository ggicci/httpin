// https://ggicci.github.io/httpin/advanced/upload-files

package httpin

import (
	"fmt"
	"mime/multipart"
	"reflect"
)

func init() {
	RegisterTypeDecoder(reflect.TypeOf(File{}), FileTypeDecoderFunc(decodeFile))
}

type File struct {
	multipart.File
	Header *multipart.FileHeader
	Valid  bool
}

func decodeFile(fileHeader *multipart.FileHeader) (interface{}, error) {
	var inFile File
	if fileHeader == nil {
		return inFile, ErrNilFile
	}

	inFile.Header = fileHeader
	file, err := fileHeader.Open()
	if err != nil {
		return inFile, fmt.Errorf("open file: %w", err)
	}

	inFile.File = file
	inFile.Valid = true
	return inFile, nil
}
