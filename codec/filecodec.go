package codec

import (
	"errors"
	"io"
	"mime/multipart"
	"net/textproto"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

// FileHeader is the interface that groups the methods of multipart.FileHeader.
type FileHeader interface {
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
	UnmarshalFile(FileHeader) error
}

type FileCodec interface {
	FileMarshaler
	FileUnmarshaler
}

func NewFileable(rv reflect.Value) (FileCodec, error) {
	if IsPatchField(rv.Type()) {
		return NewFileablePatchFieldWrapper(rv)
	}

	return newFileable(rv)
}

func newFileable(rv reflect.Value) (FileCodec, error) {
	rv, err := getPointer(rv)
	if err != nil {
		return nil, err
	}

	if rv.Type().Implements(fileCodecType) && rv.CanInterface() {
		return rv.Interface().(FileCodec), nil
	}
	return nil, errors.New("unsupported file type")
}

type FileablePatchFieldWrapper struct {
	Value            reflect.Value // of patch.Field[T]
	internalFileable FileCodec
}

func NewFileablePatchFieldWrapper(rv reflect.Value) (*FileablePatchFieldWrapper, error) {
	fileable, err := NewFileable(rv.FieldByName("Value"))
	if err != nil {
		return nil, err
	} else {
		return &FileablePatchFieldWrapper{
			Value:            rv,
			internalFileable: fileable,
		}, nil
	}
}

func (w *FileablePatchFieldWrapper) Filename() string {
	return w.internalFileable.Filename()
}

func (w *FileablePatchFieldWrapper) MarshalFile() (io.ReadCloser, error) {
	return w.internalFileable.MarshalFile()
}

func (w *FileablePatchFieldWrapper) UnmarshalFile(fh FileHeader) error {
	if err := w.internalFileable.UnmarshalFile(fh); err != nil {
		return err
	} else {
		w.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}

var fileCodecType = internal.TypeOf[FileCodec]()

// multipartFileHeader is the adaptor of multipart.FileHeader.
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

func ToFileHeaderList(fhs []*multipart.FileHeader) []FileHeader {
	result := make([]FileHeader, len(fhs))
	for i, fh := range fhs {
		result[i] = &multipartFileHeader{fh}
	}
	return result
}
