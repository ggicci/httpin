// https://ggicci.github.io/httpin/advanced/upload-files

package core

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
)

func init() {
	RegisterFileCoder[*File]()
}

// File is the builtin type of httpin to manupulate file uploads. On the server
// side, it is used to represent a file in a multipart/form-data request. On the
// client side, it is used to represent a file to be uploaded.
type File struct {
	FileHeader
	uploadFilename string
	uploadReader   io.ReadCloser
}

// UploadFile is a helper function to create a File instance from a file path.
// It is useful when you want to upload a file from the local file system.
func UploadFile(filename string) *File {
	return &File{uploadFilename: filename}
}

// UploadStream is a helper function to create a File instance from a io.Reader. It
// is useful when you want to upload a file from a stream.
func UploadStream(contentReader io.ReadCloser) *File {
	return &File{uploadReader: contentReader}
}

// Filename returns the filename of the file. On the server side, it returns the
// filename of the file in the multipart/form-data request. On the client side, it
// returns the filename of the file to be uploaded.
func (f *File) Filename() string {
	if f.IsUpload() {
		return f.uploadFilename
	}
	if f.FileHeader != nil {
		return f.FileHeader.Filename()
	}
	return ""
}

// MarshalFile implements FileMarshaler.
func (f *File) MarshalFile() (io.ReadCloser, error) {
	if f.IsUpload() {
		return f.OpenUploadStream()
	} else {
		return f.OpenReceiveStream()
	}
}

func (f *File) UnmarshalFile(fh FileHeader) error {
	f.FileHeader = fh
	return nil
}

// IsUpload returns true when the File instance is created for an upload purpose.
// Typically, you should use UploadFilename or UploadReader to create a File instance
// for upload.
func (f *File) IsUpload() bool {
	return f.uploadFilename != "" || f.uploadReader != nil
}

// OpenUploadStream returns a io.ReadCloser for the file to be uploaded.
// Call this method on the client side for uploading a file.
func (f *File) OpenUploadStream() (io.ReadCloser, error) {
	if f.uploadReader != nil {
		return f.uploadReader, nil
	}
	if f.uploadFilename != "" {
		return os.Open(f.uploadFilename)
	}
	return nil, errors.New("invalid upload (client): no filename or reader")
}

// OpenReceiveStream returns a io.Reader for the file in the multipart/form-data request.
// Call this method on the server side to read the file content.
func (f *File) OpenReceiveStream() (multipart.File, error) {
	if f.FileHeader == nil {
		return nil, errors.New("invalid upload (server): nil file header")
	}
	return f.FileHeader.Open()
}

func (f *File) ReadAll() ([]byte, error) {
	reader, err := f.MarshalFile()
	if err != nil {
		return nil, err
	}
	return io.ReadAll(reader)
}
