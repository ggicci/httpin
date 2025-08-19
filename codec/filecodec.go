package codec

import (
	"errors"
	"io"
	"mime/multipart"
	"net/textproto"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

// FileObject represents a file object, which is typically used to represent a file in
// a multipart/form-data request. See multipart.FileHeader for more details.
type FileObject interface {
	Filename() string
	Size() int64
	MIMEHeader() textproto.MIMEHeader
	Open() (multipart.File, error)
}

type FileMarshaler interface {
	Filename() string
	MarshalFile() (io.ReadCloser, error)
}

type FileUnmarshaler interface {
	UnmarshalFile(FileObject) error
}

type FileCodec interface {
	FileMarshaler
	FileUnmarshaler
}

func NewFileCodec(rv reflect.Value) (FileCodec, error) {
	if IsPatchField(rv.Type()) {
		return NewFileCodec4PatchField(rv)
	}

	return newFileCodec(rv)
}

func newFileCodec(rv reflect.Value) (FileCodec, error) {
	rv, err := getPointer(rv)
	if err != nil {
		return nil, err
	}

	if rv.Type().Implements(fileCodecType) && rv.CanInterface() {
		return rv.Interface().(FileCodec), nil
	}
	return nil, errors.New("unsupported file type")
}

var fileCodecType = internal.TypeOf[FileCodec]()

func MultipartFileHeaderAsFileObject(fh *multipart.FileHeader) FileObject {
	if fh == nil {
		return nil
	}
	return &multipartFileHeader{fh}
}

// multipartFileHeader wraps multipart.FileHeader into a FileObject.
type multipartFileHeader struct{ *multipart.FileHeader }

func (h *multipartFileHeader) Filename() string {
	return h.FileHeader.Filename
}

func (h *multipartFileHeader) Size() int64 {
	return h.FileHeader.Size
}

func (h *multipartFileHeader) MIMEHeader() textproto.MIMEHeader {
	return h.FileHeader.Header
}

func (h *multipartFileHeader) Open() (multipart.File, error) {
	return h.FileHeader.Open()
}
