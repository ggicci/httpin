package httpin

import (
	"bytes"
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

func newMultipartFormWriterFromMap(m map[string]interface{}) (body *bytes.Buffer, writer *multipart.Writer) {
	body = new(bytes.Buffer)
	writer = multipart.NewWriter(body)

	appendValue := func(key, value string) {
		fieldWriter, _ := writer.CreateFormField(key)
		fieldWriter.Write([]byte(value))
	}
	appendFile := func(key string, value []byte) {
		fileWriter, _ := writer.CreateFormFile(key, key+".txt")
		fileWriter.Write(value)
	}

	for k, v := range m {
		switch cv := v.(type) {
		case string:
			appendValue(k, cv)
		case []byte:
			appendFile(k, cv)
		case []string:
			for _, sv := range cv {
				appendValue(k, sv)
			}
		case [][]byte:
			for _, bv := range cv {
				appendFile(k, bv)
			}
		default:
			panic("invalid type")
		}
	}
	_ = writer.Close() // error ignored
	return
}

func newMultipartFormRequestFromMap(m map[string]interface{}) *http.Request {
	body, writer := newMultipartFormWriterFromMap(m)
	r, _ := http.NewRequest("POST", "/", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	return r
}

func assertFile(t *testing.T, gotFile File, filename string, content []byte) {
	assert.True(t, gotFile.Valid)
	assert.Equal(t, filename, gotFile.Header.Filename)
	assert.Equal(t, int64(len(content)), gotFile.Header.Size)
	uploadedContent, err := io.ReadAll(gotFile.File)
	assert.NoError(t, err)
	assert.Equal(t, content, uploadedContent)
}

func TestMultipartForm_UploadSingleFile(t *testing.T) {
	assert := assert.New(t)
	// Upload a file through multipart/form-data requests.
	var AvatarBytes = []byte("avatar image content")

	r := newMultipartFormRequestFromMap(map[string]interface{}{
		"name":   "Ggicci T'ang",
		"gender": "male",
		"avatar": AvatarBytes,
	})
	core, err := New(UpdateUserProfileInput{})
	assert.NoError(err)
	gotInput, err := core.Decode(r)
	assert.NoError(err)
	got, ok := gotInput.(*UpdateUserProfileInput)
	assert.True(ok)
	assert.Equal("Ggicci T'ang", got.Name)
	assert.Equal("male", got.Gender)
	assertFile(t, got.Avatar, "avatar.txt", AvatarBytes)
}

func TestMultipartForm_UploadSingleFile_FailOnEmpty(t *testing.T) {
	assert := assert.New(t)
	r := newMultipartFormRequestFromMap(map[string]interface{}{
		"name": "Ggicci T'ang",
		// No files uploaded should cause server to fail.
	})
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

func breakMultipartFormBoundary(body *bytes.Buffer) *bytes.Buffer {
	raw := body.Bytes()
	var brokenBody = bytes.NewBuffer(raw[:len(raw)-10])
	brokenBody.Write([]byte("xxx")) // break the boundary
	return brokenBody
}

func TestMultipartForm_UploadSingleFile_FailOnBrokenBoundaries(t *testing.T) {
	assert := assert.New(t)
	// Broken boundaries should cause server to fail.
	var AvatarBytes = []byte("avatar image content")
	body, writer := newMultipartFormWriterFromMap(map[string]interface{}{
		"avatar": AvatarBytes,
	})

	r, _ := http.NewRequest("POST", "/", breakMultipartFormBoundary(body))
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
	title := "feature-request: integrate with open-telemetry"
	var attachments = [][]byte{
		[]byte("attachment #1"),
		[]byte("attachment #2"),
		[]byte("attachment #3"),
	}

	r := newMultipartFormRequestFromMap(map[string]interface{}{
		"title":      title,
		"attachment": attachments,
	})
	core, err := New(UpdateGitHubIssueInput{})
	assert.NoError(err)
	gotInput, err := core.Decode(r)
	assert.NoError(err)
	got, ok := gotInput.(*UpdateGitHubIssueInput)
	assert.True(ok)
	assert.Equal(title, got.Title)
	assert.Len(got.Attachments, len(attachments))
	for i, att := range got.Attachments {
		assertFile(t, att, "attachment.txt", attachments[i])
	}
}
