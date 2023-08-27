package httpin

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UpdateUserProfileInput struct {
	Name   string `in:"form=name"`
	Gender string `in:"form=gender"`
	Avatar File   `in:"form=avatar"`
}

type UpdateGitHubIssueInput struct {
	Title       string `in:"form=title"`
	Attachments []File `in:"form=attachment"`
}

func TestMultipartForm_DecodeFile_FailOnNilFileHeader(t *testing.T) {
	gotInput, err := decodeFile(nil)
	assert.ErrorIs(t, err, ErrNilFile)
	got, ok := gotInput.(File)
	assert.True(t, ok)
	assert.False(t, got.Valid)
}

func TestMultipartForm_DecodeFile_FailOnBrokenFileHeader(t *testing.T) {
	fileHeader := &multipart.FileHeader{
		Filename: "avatar.png",
		Size:     10,
	}
	gotInput, err := decodeFile(fileHeader)
	assert.Error(t, err)
	got, ok := gotInput.(File)
	assert.True(t, ok)
	assert.False(t, got.Valid)
}

func TestMultipartForm_UploadSingleFile(t *testing.T) {
	assert := assert.New(t)
	// Upload a file through multipart/form-data requests.
	var AvatarBytes = []byte("avatar image content")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	nameFieldWriter, err := writer.CreateFormField("name")
	assert.NoError(err)
	nameFieldWriter.Write([]byte("Ggicci T'ang"))

	genderFieldWriter, err := writer.CreateFormField("gender")
	assert.NoError(err)
	genderFieldWriter.Write([]byte("male"))

	avatarFileWriter, err := writer.CreateFormFile("avatar", "avatar.png")
	assert.NoError(err)
	_, err = avatarFileWriter.Write(AvatarBytes)
	assert.NoError(err)

	_ = writer.Close() // error ignored

	r, _ := http.NewRequest("POST", "/", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	core, err := New(UpdateUserProfileInput{})
	assert.NoError(err)
	gotInput, err := core.Decode(r)
	assert.NoError(err)
	got, ok := gotInput.(*UpdateUserProfileInput)
	assert.True(ok)
	assert.Equal("Ggicci T'ang", got.Name)
	assert.Equal("male", got.Gender)
	assert.True(got.Avatar.Valid)
	assert.Equal("avatar.png", got.Avatar.Header.Filename)
	assert.Equal(int64(len(AvatarBytes)), got.Avatar.Header.Size)
	uploadedContent, err := io.ReadAll(got.Avatar.File)
	assert.NoError(err)
	assert.Equal(AvatarBytes, uploadedContent)
}

func TestMultipartForm_UploadSingleFile_FailOnEmpty(t *testing.T) {
	assert := assert.New(t)
	// No files uploaded should cause server to fail.
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	nameFieldWriter, err := writer.CreateFormField("name")
	assert.NoError(err)
	nameFieldWriter.Write([]byte("Ggicci T'ang"))

	_ = writer.Close() // error ignored

	r, _ := http.NewRequest("POST", "/", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	core, err := New(UpdateUserProfileInput{})
	assert.NoError(err)
	gotInput, err := core.Decode(r)
	assert.NoError(err)
	got, ok := gotInput.(*UpdateUserProfileInput)
	assert.True(ok)

	assert.Equal("Ggicci T'ang", got.Name)
	assert.False(got.Avatar.Valid)
	assert.Nil(got.Avatar.File)
	assert.Nil(got.Avatar.Header)
}

func TestMultipartForm_UploadSingleFile_FailOnBrokenBoundaries(t *testing.T) {
	assert := assert.New(t)
	// Broken boundaries should cause server to fail.
	var AvatarBytes = []byte("avatar image content")

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	avatarFileWriter, err := writer.CreateFormFile("avatar", "avatar.png")
	assert.NoError(err)
	_, err = avatarFileWriter.Write(AvatarBytes)
	assert.NoError(err)
	writer.Close() // error ignored

	raw := body.Bytes()
	var brokenBody = bytes.NewBuffer(raw[:len(raw)-10])
	brokenBody.Write([]byte("xxx")) // break the boundary

	r, _ := http.NewRequest("POST", "/", brokenBody)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	core, err := New(UpdateUserProfileInput{})
	assert.NoError(err)

	gotInput, err := core.Decode(r)
	assert.Nil(gotInput)
	assert.Error(err)
}

func TestMultipartForm_UploadMultiFiles(t *testing.T) {
	assert := assert.New(t)
	// Upload multiple files at a time.
	var attachments = [][]byte{
		[]byte("attachment #1"),
		[]byte("attachment #2"),
		[]byte("attachment #3"),
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	title := "feature-request: integrate with open-telemetry"
	titleFieldWriter, err := writer.CreateFormField("title")
	assert.NoError(err)
	titleFieldWriter.Write([]byte(title))

	for i, attContent := range attachments {
		filename := fmt.Sprintf("attachment-%d.txt", i+1)
		attachmentFileWriter, err := writer.CreateFormFile("attachment", filename)
		assert.NoError(err)
		_, err = attachmentFileWriter.Write(attContent)
		assert.NoError(err)
	}
	_ = writer.Close() // error ignored

	r, _ := http.NewRequest("POST", "/", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	core, err := New(UpdateGitHubIssueInput{})
	assert.NoError(err)
	gotInput, err := core.Decode(r)
	assert.NoError(err)
	got, ok := gotInput.(*UpdateGitHubIssueInput)
	assert.True(ok)
	assert.Equal(title, got.Title)
	assert.Len(got.Attachments, len(attachments))
	for i, att := range got.Attachments {
		assert.True(att.Valid)
		assert.Equal(fmt.Sprintf("attachment-%d.txt", i+1), att.Header.Filename)
		assert.Equal(int64(len(attachments[i])), att.Header.Size)
		uploadedContent, err := io.ReadAll(att.File)
		assert.NoError(err)
		assert.Equal(attachments[i], uploadedContent)
	}
}
