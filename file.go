// https://ggicci.github.io/httpin/advanced/upload-files

package httpin

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
)

var fileTypes = make(map[reflect.Type]fileDecoderAdaptor)

func init() {
	RegisterFileType[*File](FileDecoderFunc[*File](decodeFile))
}

// FileEncoder is the interface implemented by types that can represent a file upload.
type FileEncoder interface {
	// Encode returns the filename of the file and a io.ReadCloser for the file content.
	Encode() (string, io.ReadCloser, error)
}

// RegisterFileType registers a FileEncodeDecoder for type T. Which marks the type T as
// a file type. When httpin encounters a field of type T, it will treat it as a file
// upload.
//
//	func init() {
//	    RegisterFileType[MyFile](&myFileEncodeDecoder{})
//	}
func RegisterFileType[T FileEncoder](fd FileDecoder[T]) {
	typ := typeOf[T]()
	if isFileType(typ) {
		panicOnError(fmt.Errorf("duplicate file type: %v", typ))
	}
	if fd == nil {
		panicOnError(errors.New("nil decoder"))
	}
	fileTypes[typ] = adaptDecoder(typ, toAnyFileDecoder(fd)).(fileDecoderAdaptor)
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

func fileUploadBuilder(rtm *DirectiveRuntime, files []FileEncoder) error {
	rb := rtm.GetRequestBuilder()
	key := rtm.Directive.Argv[0]
	rb.setAttachment(key, files)
	rtm.MarkFieldSet(true)
	return nil
}

func toAnyFileDecoder[T FileEncoder](fd FileDecoder[T]) FileDecoder[any] {
	return FileDecoderFunc[any](func(fh *multipart.FileHeader) (any, error) {
		return fd.Decode(fh)
	})
}

func fileDecoderByType(typ reflect.Type) fileDecoderAdaptor {
	return fileTypes[typ]
}

func isFileType(typ reflect.Type) bool {
	return fileDecoderByType(typ) != nil
}
