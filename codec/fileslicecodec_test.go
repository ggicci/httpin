package codec

import (
	"reflect"
	"testing"

	"github.com/ggicci/httpin/internal/testutil"
	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

func TestFileSlicable_FromFileSlice(t *testing.T) {
	rv := reflect.New(reflect.TypeOf(MyFiles{})).Elem()
	s := rv.Addr().Interface().(*MyFiles)

	fileAvatar := testutil.CreateTempFileV2(t)
	fileAvatarPointer := testutil.CreateTempFileV2(t)
	filePatchAvatar := testutil.CreateTempFileV2(t)
	filePatchAvatarPointer := testutil.CreateTempFileV2(t)

	testAssignFileSlice(t, rv.FieldByName("Avatar"), []FileObject{
		mockFileHeader(t, fileAvatar.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("AvatarPointer"), []FileObject{
		mockFileHeader(t, fileAvatarPointer.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("Avatars"), []FileObject{
		mockFileHeader(t, fileAvatar.Filename),
		mockFileHeader(t, fileAvatarPointer.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("AvatarPointers"), []FileObject{
		mockFileHeader(t, fileAvatarPointer.Filename),
		mockFileHeader(t, fileAvatar.Filename),
	})

	testAssignFileSlice(t, rv.FieldByName("PatchAvatar"), []FileObject{
		mockFileHeader(t, filePatchAvatar.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("PatchAvatarPointer"), []FileObject{
		mockFileHeader(t, filePatchAvatarPointer.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("PatchAvatars"), []FileObject{
		mockFileHeader(t, fileAvatar.Filename),
		mockFileHeader(t, filePatchAvatar.Filename),
		mockFileHeader(t, filePatchAvatarPointer.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("PatchAvatarPointers"), []FileObject{
		mockFileHeader(t, fileAvatar.Filename),
		mockFileHeader(t, fileAvatarPointer.Filename),
		mockFileHeader(t, filePatchAvatar.Filename),
		mockFileHeader(t, filePatchAvatarPointer.Filename),
	})

	validateFile(t, fileAvatar, &s.Avatar)
	validateFile(t, fileAvatarPointer, s.AvatarPointer)

	assert.Len(t, s.Avatars, 2)
	validateFile(t, fileAvatar, &s.Avatars[0])
	validateFile(t, fileAvatarPointer, &s.Avatars[1])

	assert.Len(t, s.AvatarPointers, 2)
	validateFile(t, fileAvatarPointer, s.AvatarPointers[0])
	validateFile(t, fileAvatar, s.AvatarPointers[1])

	assert.True(t, s.PatchAvatar.Valid)
	validateFile(t, filePatchAvatar, &s.PatchAvatar.Value)

	assert.True(t, s.PatchAvatarPointer.Valid)
	validateFile(t, filePatchAvatarPointer, s.PatchAvatarPointer.Value)

	assert.True(t, s.PatchAvatars.Valid)
	assert.Len(t, s.PatchAvatars.Value, 3)
	validateFile(t, fileAvatar, &s.PatchAvatars.Value[0])
	validateFile(t, filePatchAvatar, &s.PatchAvatars.Value[1])
	validateFile(t, filePatchAvatarPointer, &s.PatchAvatars.Value[2])

	assert.True(t, s.PatchAvatarPointers.Valid)
	assert.Len(t, s.PatchAvatarPointers.Value, 4)
	validateFile(t, fileAvatar, s.PatchAvatarPointers.Value[0])
	validateFile(t, fileAvatarPointer, s.PatchAvatarPointers.Value[1])
	validateFile(t, filePatchAvatar, s.PatchAvatarPointers.Value[2])
	validateFile(t, filePatchAvatarPointer, s.PatchAvatarPointers.Value[3])
}

func TestFileSlicable_ToFileSlice(t *testing.T) {
	fileAvatar := testutil.CreateTempFileV2(t)
	fileAvatarPointer := testutil.CreateTempFileV2(t)
	filePatchAvatar := testutil.CreateTempFileV2(t)
	filePatchAvatarPointer := testutil.CreateTempFileV2(t)

	var s = &MyFiles{
		Avatar:         *UploadFile(fileAvatar.Filename),
		AvatarPointer:  UploadFile(fileAvatarPointer.Filename),
		Avatars:        []File{*UploadFile(fileAvatar.Filename), *UploadFile(fileAvatarPointer.Filename)},
		AvatarPointers: []*File{UploadFile(fileAvatarPointer.Filename), UploadFile(fileAvatar.Filename)},
		PatchAvatar:    patch.Field[File]{Value: *UploadFile(filePatchAvatar.Filename), Valid: true},
		PatchAvatarPointer: patch.Field[*File]{
			Value: UploadFile(filePatchAvatarPointer.Filename),
			Valid: true,
		},
		PatchAvatars: patch.Field[[]File]{
			Value: []File{
				*UploadFile(fileAvatar.Filename),
				*UploadFile(filePatchAvatar.Filename),
				*UploadFile(filePatchAvatarPointer.Filename),
			},
			Valid: true,
		},
		PatchAvatarPointers: patch.Field[[]*File]{
			Value: []*File{
				UploadFile(fileAvatar.Filename),
				UploadFile(fileAvatarPointer.Filename),
				UploadFile(filePatchAvatar.Filename),
				UploadFile(filePatchAvatarPointer.Filename),
			},
			Valid: true,
		},
	}

	rv := reflect.ValueOf(s).Elem()
	validateFileList(t, []*testutil.NamedTempFile{fileAvatar}, testGetFileSlice(t, rv.FieldByName("Avatar")))
	validateFileList(t, []*testutil.NamedTempFile{fileAvatarPointer}, testGetFileSlice(t, rv.FieldByName("AvatarPointer")))
	validateFileList(t, []*testutil.NamedTempFile{fileAvatar, fileAvatarPointer}, testGetFileSlice(t, rv.FieldByName("Avatars")))
	validateFileList(t, []*testutil.NamedTempFile{fileAvatarPointer, fileAvatar}, testGetFileSlice(t, rv.FieldByName("AvatarPointers")))
	validateFileList(t, []*testutil.NamedTempFile{filePatchAvatar}, testGetFileSlice(t, rv.FieldByName("PatchAvatar")))
	validateFileList(t, []*testutil.NamedTempFile{filePatchAvatarPointer}, testGetFileSlice(t, rv.FieldByName("PatchAvatarPointer")))
	validateFileList(t, []*testutil.NamedTempFile{fileAvatar, filePatchAvatar, filePatchAvatarPointer}, testGetFileSlice(t, rv.FieldByName("PatchAvatars")))
	validateFileList(t, []*testutil.NamedTempFile{fileAvatar, fileAvatarPointer, filePatchAvatar, filePatchAvatarPointer}, testGetFileSlice(t, rv.FieldByName("PatchAvatarPointers")))
}

func testAssignFileSlice(t *testing.T, rv reflect.Value, files []FileObject) {
	fs, err := NewFileSliceCodec(rv)
	assert.NoError(t, err)
	assert.NoError(t, fs.FromFileSlice(files))
}

func testGetFileSlice(t *testing.T, rv reflect.Value) []FileMarshaler {
	fs, err := NewFileSliceCodec(rv)
	assert.NoError(t, err)
	files, err := fs.ToFileSlice()
	assert.NoError(t, err)
	return files
}

func validateFileList(t *testing.T, expected []*testutil.NamedTempFile, actual []FileMarshaler) {
	assert.Len(t, actual, len(expected))
	for i, file := range expected {
		validateFile(t, file, actual[i])
	}
}
