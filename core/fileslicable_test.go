package core

import (
	"reflect"
	"testing"

	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

func TestFileSlicable_FromFileSlice(t *testing.T) {
	rv := reflect.New(reflect.TypeOf(MyFiles{})).Elem()
	s := rv.Addr().Interface().(*MyFiles)

	fileAvatar := createTempFileV2(t)
	fileAvatarPointer := createTempFileV2(t)
	filePatchAvatar := createTempFileV2(t)
	filePatchAvatarPointer := createTempFileV2(t)

	testAssignFileSlice(t, rv.FieldByName("Avatar"), []FileHeader{
		mockFileHeader(t, fileAvatar.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("AvatarPointer"), []FileHeader{
		mockFileHeader(t, fileAvatarPointer.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("Avatars"), []FileHeader{
		mockFileHeader(t, fileAvatar.Filename),
		mockFileHeader(t, fileAvatarPointer.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("AvatarPointers"), []FileHeader{
		mockFileHeader(t, fileAvatarPointer.Filename),
		mockFileHeader(t, fileAvatar.Filename),
	})

	testAssignFileSlice(t, rv.FieldByName("PatchAvatar"), []FileHeader{
		mockFileHeader(t, filePatchAvatar.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("PatchAvatarPointer"), []FileHeader{
		mockFileHeader(t, filePatchAvatarPointer.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("PatchAvatars"), []FileHeader{
		mockFileHeader(t, fileAvatar.Filename),
		mockFileHeader(t, filePatchAvatar.Filename),
		mockFileHeader(t, filePatchAvatarPointer.Filename),
	})
	testAssignFileSlice(t, rv.FieldByName("PatchAvatarPointers"), []FileHeader{
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
	fileAvatar := createTempFileV2(t)
	fileAvatarPointer := createTempFileV2(t)
	filePatchAvatar := createTempFileV2(t)
	filePatchAvatarPointer := createTempFileV2(t)

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
	validateFileList(t, []*tempFile{fileAvatar}, testGetFileSlice(t, rv.FieldByName("Avatar")))
	validateFileList(t, []*tempFile{fileAvatarPointer}, testGetFileSlice(t, rv.FieldByName("AvatarPointer")))
	validateFileList(t, []*tempFile{fileAvatar, fileAvatarPointer}, testGetFileSlice(t, rv.FieldByName("Avatars")))
	validateFileList(t, []*tempFile{fileAvatarPointer, fileAvatar}, testGetFileSlice(t, rv.FieldByName("AvatarPointers")))
	validateFileList(t, []*tempFile{filePatchAvatar}, testGetFileSlice(t, rv.FieldByName("PatchAvatar")))
	validateFileList(t, []*tempFile{filePatchAvatarPointer}, testGetFileSlice(t, rv.FieldByName("PatchAvatarPointer")))
	validateFileList(t, []*tempFile{fileAvatar, filePatchAvatar, filePatchAvatarPointer}, testGetFileSlice(t, rv.FieldByName("PatchAvatars")))
	validateFileList(t, []*tempFile{fileAvatar, fileAvatarPointer, filePatchAvatar, filePatchAvatarPointer}, testGetFileSlice(t, rv.FieldByName("PatchAvatarPointers")))
}

func testAssignFileSlice(t *testing.T, rv reflect.Value, files []FileHeader) {
	fs, err := NewFileSlicable(rv)
	assert.NoError(t, err)
	assert.NoError(t, fs.FromFileSlice(files))
}

func testGetFileSlice(t *testing.T, rv reflect.Value) []FileMarshaler {
	fs, err := NewFileSlicable(rv)
	assert.NoError(t, err)
	files, err := fs.ToFileSlice()
	assert.NoError(t, err)
	return files
}

func validateFileList(t *testing.T, expected []*tempFile, actual []FileMarshaler) {
	assert.Len(t, actual, len(expected))
	for i, file := range expected {
		validateFile(t, file, actual[i])
	}
}
