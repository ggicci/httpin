// https://ggicci.github.io/httpin/advanced/upload-files

package core

import (
	"io"

	"github.com/ggicci/httpin/codec"
)

// File is the builtin type of httpin to manupulate file uploads. On the server
// side, it is used to represent a file in a multipart/form-data request. On the
// client side, it is used to represent a file to be uploaded.
type File = codec.File

func init() {
	RegisterFileCodec[*File]()
}

// UploadFile is a helper function to create a File instance from a file path.
// It is useful when you want to upload a file from the local file system.
func UploadFile(filename string) *File {
	return codec.UploadFile(filename)
}

// UploadStream is a helper function to create a File instance from a io.Reader. It
// is useful when you want to upload a file from a stream.
func UploadStream(contentReader io.ReadCloser) *File {
	return codec.UploadStream(contentReader)
}
