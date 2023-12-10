package core

import (
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"reflect"
	"testing"

	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

type MyFiles struct {
	Avatar         File
	AvatarPointer  *File
	Avatars        []File
	AvatarPointers []*File

	PatchAvatar         patch.Field[File]
	PatchAvatarPointer  patch.Field[*File]
	PatchAvatars        patch.Field[[]File]
	PatchAvatarPointers patch.Field[[]*File]
}

func TestFileable_UnmarshalFile(t *testing.T) {
	rv := reflect.New(reflect.TypeOf(MyFiles{})).Elem()
	s := rv.Addr().Interface().(*MyFiles)

	fileAvatar := testAssignFile(t, rv.FieldByName("Avatar"))
	fileAvatarPointer := testAssignFile(t, rv.FieldByName("AvatarPointer"))
	testNewFileableErrUnsupported(t, rv.FieldByName("Avatars"))
	testNewFileableErrUnsupported(t, rv.FieldByName("AvatarPointers"))

	filePatchAvatar := testAssignFile(t, rv.FieldByName("PatchAvatar"))
	filePatchAvatarPointer := testAssignFile(t, rv.FieldByName("PatchAvatarPointer"))
	testNewFileableErrUnsupported(t, rv.FieldByName("PatchAvatars"))
	testNewFileableErrUnsupported(t, rv.FieldByName("PatchAvatarPointers"))

	validateFile(t, fileAvatar, &s.Avatar)
	validateFile(t, fileAvatarPointer, s.AvatarPointer)

	assert.True(t, s.PatchAvatar.Valid)
	validateFile(t, filePatchAvatar, &s.PatchAvatar.Value)
	assert.True(t, s.PatchAvatarPointer.Valid)
	validateFile(t, filePatchAvatarPointer, s.PatchAvatarPointer.Value)
}

func TestFileable_MarshalFile(t *testing.T) {
	fileAvatar := createTempFileV2(t)
	fileAvatarPointer := createTempFileV2(t)
	filePatchAvatar := createTempFileV2(t)
	filePatchAvatarPointer := createTempFileV2(t)

	var s = &MyFiles{
		Avatar:         *UploadFile(fileAvatar.Filename),
		AvatarPointer:  UploadFile(fileAvatarPointer.Filename),
		Avatars:        []File{*UploadFile(fileAvatar.Filename)},
		AvatarPointers: []*File{UploadFile(fileAvatarPointer.Filename)},

		PatchAvatar:         patch.Field[File]{Value: *UploadFile(filePatchAvatar.Filename), Valid: true},
		PatchAvatarPointer:  patch.Field[*File]{Value: UploadFile(filePatchAvatarPointer.Filename), Valid: true},
		PatchAvatars:        patch.Field[[]File]{Value: []File{*UploadFile(filePatchAvatar.Filename)}, Valid: true},
		PatchAvatarPointers: patch.Field[[]*File]{Value: []*File{UploadFile(filePatchAvatarPointer.Filename)}, Valid: true},
	}

	rv := reflect.ValueOf(s).Elem()

	assert.Equal(t, fileAvatar.Content, testGetFile(t, rv.FieldByName("Avatar")))
	assert.Equal(t, fileAvatarPointer.Content, testGetFile(t, rv.FieldByName("AvatarPointer")))
	testNewFileableErrUnsupported(t, rv.FieldByName("Avatars"))
	testNewFileableErrUnsupported(t, rv.FieldByName("AvatarPointers"))

	assert.Equal(t, filePatchAvatar.Content, testGetFile(t, rv.FieldByName("PatchAvatar")))
	assert.Equal(t, filePatchAvatarPointer.Content, testGetFile(t, rv.FieldByName("PatchAvatarPointer")))
	testNewFileableErrUnsupported(t, rv.FieldByName("PatchAvatars"))
	testNewFileableErrUnsupported(t, rv.FieldByName("PatchAvatarPointers"))
}

func testNewFileableErrUnsupported(t *testing.T, rv reflect.Value) {
	fileable, err := NewFileable(rv)
	assert.ErrorContains(t, err, "unsupported file type")
	assert.Nil(t, fileable)
}

func validateFile(t *testing.T, expected *tempFile, actual FileMarshaler) {
	assert.Equal(t, expected.Filename, actual.Filename())
	reader, err := actual.MarshalFile()
	assert.NoError(t, err)
	content, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, expected.Content, content)
}

func testAssignFile(t *testing.T, rv reflect.Value) *tempFile {
	fileable, err := NewFileable(rv)
	assert.NoError(t, err)
	file := createTempFileV2(t)
	assert.NoError(t, fileable.UnmarshalFile(mockFileHeader(t, file.Filename)))
	return file
}

func testGetFile(t *testing.T, rv reflect.Value) []byte {
	file, err := NewFileable(rv)
	assert.NoError(t, err)
	reader, err := file.MarshalFile()
	assert.NoError(t, err)
	content, err := io.ReadAll(reader)
	assert.NoError(t, err)
	return content
}

type dummyFileHeader struct {
	file *os.File
}

func mockFileHeader(t *testing.T, filename string) FileHeader {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return &dummyFileHeader{
		file: file,
	}
}

func (f *dummyFileHeader) Filename() string {
	return f.file.Name()
}

func (f *dummyFileHeader) Size() int64 {
	stat, err := f.file.Stat()
	if err != nil {
		panic(err)
	}
	return stat.Size()
}

func (f *dummyFileHeader) MIMEHeader() textproto.MIMEHeader {
	return textproto.MIMEHeader{}
}

func (f *dummyFileHeader) Open() (multipart.File, error) {
	return f.file, nil
}
