// https://ggicci.github.io/httpin/advanced/upload-files

package httpin

import (
	"mime/multipart"
)

func init() {
	registerTypeDecoderTo[File](builtinDecoders, DecoderFunc[*multipart.FileHeader](decodeFile), false)
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
		return inFile, err
	}

	inFile.File = file
	inFile.Valid = true
	return inFile, nil
}
