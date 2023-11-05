// https://ggicci.github.io/httpin/advanced/upload-files

package internal

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
)

// FileEncoder is the interface implemented by types that can represent a file upload.
type FileEncoder interface {
	// Encode returns the filename of the file and a io.ReadCloser for the file content.
	Encode() (string, io.ReadCloser, error)
}

// File is the builtin type of httpin to manupulate file uploads. On the server
// side, it is used to represent a file in a multipart/form-data request. On the
// client side, it is used to represent a file to be uploaded.
type File struct {
	Header         *multipart.FileHeader
	uploadFilename string
	uploadReader   io.ReadCloser
}

// UploadWithFilename is a helper function to create a File instance from a file path.
// It is useful when you want to upload a file from the local file system.
func UploadWithFilename(filename string) *File {
	return &File{uploadFilename: filename}
}

// UploadWithReader is a helper function to create a File instance from a io.Reader. It
// is useful when you want to upload a file from a stream.
func UploadWithReader(contentReader io.ReadCloser) *File {
	return &File{uploadReader: contentReader}
}

// Encode implements FileEncoder.
func (f File) Encode() (string, io.ReadCloser, error) {
	uploadReader, err := f.OpenUploadStream()
	if err != nil {
		return "", nil, err
	}
	return f.Filename(), uploadReader, nil
}

// Filename returns the filename of the file. On the server side, it returns the
// filename of the file in the multipart/form-data request. On the client side, it
// returns the filename of the file to be uploaded.
func (f *File) Filename() string {
	if f.IsUpload() {
		return f.uploadFilename
	}
	return f.Header.Filename
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
	if f.Header == nil {
		return nil, errors.New("invalid upload (server): nil file header")
	}
	return f.Header.Open()
}

func decodeFile(fileHeader *multipart.FileHeader) (*File, error) {
	if fileHeader == nil {
		return nil, errors.New("nil file header")
	}
	return &File{Header: fileHeader}, nil
}

func ToAnyFileDecoder[T FileEncoder](fd FileDecoder[T]) FileDecoder[any] {
	if fd == nil {
		return nil
	}
	return FileDecoderFunc[any](func(fh *multipart.FileHeader) (any, error) {
		return fd.Decode(fh)
	})
}

func toFileEncoders(value reflect.Value, kind TypeKind) ([]FileEncoder, error) {
	if IsNil(value) {
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
	if IsNil(value) {
		return errors.New("nil file encoder")
	}
	return nil
}
