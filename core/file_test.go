package core

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ggicci/httpin/internal"
	"github.com/stretchr/testify/assert"
)

type BadFile struct{}

func (bf *BadFile) Filename() string {
	return ""
}

func (bf *BadFile) MarshalFile() (io.ReadCloser, error) {
	return nil, nil
}

var errBadFile = errors.New("bad file")

func (bf *BadFile) UnmarshalFile(fh internal.FileHeader) error {
	return errBadFile
}

type UpdateUserProfileInput struct {
	Name   string `in:"form=name"`
	Gender string `in:"form=gender"`
	Avatar *File  `in:"form=avatar"`
}

type UpdateGitHubIssueInput struct {
	Title       string  `in:"form=title"`
	Attachments []*File `in:"form=attachment"`
}

func TestRegisterFileCoder(t *testing.T) {
	RegisterFileCodec[*File]()
	RegisterFileCodec[*File]() // register twice is ok
}

func TestFile_OpenUploadStream_FailOnInvalidUpload(t *testing.T) {
	file := &File{} // invalid upload, must use UploadWithFilename or UploadWithReader
	_, err := file.OpenUploadStream()
	assert.ErrorContains(t, err, "invalid upload (client)")
}

func TestFile_ReceiveStream_FailOnInvalidUpload(t *testing.T) {
	file := &File{}
	_, err := file.OpenReceiveStream()
	assert.ErrorContains(t, err, "invalid upload (server)")
}

func TestMultipartForm_UploadSingleFile(t *testing.T) {
	assert := assert.New(t)
	// Upload a file through multipart/form-data requests.
	var AvatarBytes = []byte("avatar image content")

	r := newMultipartFormRequestFromMap(map[string]any{
		"name":   "Ggicci T'ang",
		"gender": "male",
		"avatar": AvatarBytes,
	})
	co, err := New(UpdateUserProfileInput{})
	assert.NoError(err)
	gotInput, err := co.Decode(r)
	assert.NoError(err)
	got, ok := gotInput.(*UpdateUserProfileInput)
	assert.True(ok)
	assert.Equal("Ggicci T'ang", got.Name)
	assert.Equal("male", got.Gender)
	assertDecodedFile(t, got.Avatar, "avatar.txt", AvatarBytes)
}

func TestMultipartForm_UploadSingleFile_FailOnEmpty(t *testing.T) {
	assert := assert.New(t)
	r := newMultipartFormRequestFromMap(map[string]any{
		"name": "Ggicci T'ang",
		// No files uploaded should cause server to fail.
	})
	co, err := New(UpdateUserProfileInput{})
	assert.NoError(err)
	gotInput, err := co.Decode(r)
	assert.NoError(err)
	got, ok := gotInput.(*UpdateUserProfileInput)
	assert.True(ok)

	assert.Equal("Ggicci T'ang", got.Name)
	assert.Nil(got.Avatar)
}

func TestMultipartForm_UploadSingleFile_FailOnBrokenBoundaries(t *testing.T) {
	// Broken boundaries will break when parsing multipart/form-data requests.
	// Which means it will fail before stepping into the Resolve process.
	assert := assert.New(t)
	// Broken boundaries should cause server to fail.
	var AvatarBytes = []byte("avatar image content")
	body, writer := newMultipartFormWriterFromMap(map[string]any{
		"avatar": AvatarBytes,
	})

	r, _ := http.NewRequest("POST", "/", breakMultipartFormBoundary(body))
	r.Header.Set("Content-Type", writer.FormDataContentType())
	co, err := New(UpdateUserProfileInput{})
	assert.NoError(err)

	gotInput, err := co.Decode(r)
	assert.Nil(gotInput)
	assert.Error(err)
}

func TestMultipartForm_UploadSingleFile_FailOnDecodeError(t *testing.T) {
	RegisterFileCodec[*BadFile]()

	type AccountUpdate struct {
		Username string   `in:"form=username"`
		Avatar   *BadFile `in:"form=avatar"`
	}

	assert := assert.New(t)
	r := newMultipartFormRequestFromMap(map[string]any{
		"username": "ggicci",
		"avatar":   []byte("avatar image content"),
	})
	co, err := New(AccountUpdate{})
	assert.NoError(err)
	file, err := co.Decode(r)
	assert.Nil(file)
	assert.ErrorIs(err, errBadFile)

	removeFileType[*BadFile]()
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

	r := newMultipartFormRequestFromMap(map[string]any{
		"title":      title,
		"attachment": attachments,
	})
	co, err := New(UpdateGitHubIssueInput{})
	assert.NoError(err)
	gotInput, err := co.Decode(r)
	assert.NoError(err)
	got, ok := gotInput.(*UpdateGitHubIssueInput)
	assert.True(ok)
	assert.Equal(title, got.Title)
	assert.Len(got.Attachments, len(attachments))
	for i, att := range got.Attachments {
		assertDecodedFile(t, att, "attachment.txt", attachments[i])
	}
}

func TestMultipartFormEncode_UploadFilename(t *testing.T) {
	type Post struct {
		Username string  `in:"form=username"`
		Main     *File   `in:"form=main"`
		Pictures []*File `in:"form=pictures"`
	}

	// Client side: upload files (encode).
	mainFilename := createTempFile(t, []byte("main content"))
	pic1Filename := createTempFile(t, []byte("pic1 content"))
	pic2Filename := createTempFile(t, []byte("pic2 content"))

	payload := &Post{
		Username: "ggicci",
		Main:     UploadFile(mainFilename),
		Pictures: []*File{
			UploadFile(pic1Filename),
			UploadFile(pic2Filename),
		},
	}
	co, err := New(Post{})
	assert.NoError(t, err)
	req, err := co.NewRequest("POST", "/post", payload)
	assert.NoError(t, err)

	// Server side: receive files (decode).
	gotValue, err := co.Decode(req)
	assert.NoError(t, err)
	got, ok := gotValue.(*Post)
	assert.True(t, ok)
	assert.Equal(t, "ggicci", got.Username)
	assertDecodedFile(t, got.Main, filepath.Base(mainFilename), []byte("main content"))
	assert.Len(t, got.Pictures, 2)
	assertDecodedFile(t, got.Pictures[0], filepath.Base(pic1Filename), []byte("pic1 content"))
	assertDecodedFile(t, got.Pictures[1], filepath.Base(pic2Filename), []byte("pic2 content"))
}

func TestMultipartFormEncode_UploadReader(t *testing.T) {
	type Post struct {
		Username string  `in:"form=username"`
		Main     *File   `in:"form=main"`
		Pictures []*File `in:"form=pictures"`
	}

	// Client side: upload files (encode).
	mainReader := bytes.NewReader([]byte("main content"))
	pic1Reader := bytes.NewReader([]byte("pic1 content"))
	pic2Reader := bytes.NewReader([]byte("pic2 content"))

	payload := &Post{
		Username: "ggicci",
		Main:     UploadStream(io.NopCloser(mainReader)),
		Pictures: []*File{
			UploadStream(io.NopCloser(pic1Reader)),
			UploadStream(io.NopCloser(pic2Reader)),
		},
	}
	co, err := New(Post{})
	assert.NoError(t, err)
	req, err := co.NewRequest("POST", "/post", payload)
	assert.NoError(t, err)

	// Server side: receive files (decode).
	gotValue, err := co.Decode(req)
	assert.NoError(t, err)
	got, ok := gotValue.(*Post)
	assert.True(t, ok)
	assert.Equal(t, "ggicci", got.Username)
	assertDecodedFile(t, got.Main, "main_0", []byte("main content"))
	assert.Len(t, got.Pictures, 2)
	assertDecodedFile(t, got.Pictures[0], "pictures_0", []byte("pic1 content"))
	assertDecodedFile(t, got.Pictures[1], "pictures_1", []byte("pic2 content"))
}

func TestUpload_WithNilFile(t *testing.T) {
	payload := &UpdateUserProfileInput{
		Avatar: nil,
	}
	expected, _ := http.NewRequest("POST", "/post", nil)
	expectedForm := url.Values{
		"name":   {""},
		"gender": {""},
	}
	expected.Body = io.NopCloser(strings.NewReader(expectedForm.Encode()))
	expected.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	co, err := New(UpdateUserProfileInput{})
	assert.NoError(t, err)
	req, err := co.NewRequest("POST", "/post", payload)
	assert.NoError(t, err)
	assert.Equal(t, expected, req)
}

func TestUpload_WithNilMultiFile(t *testing.T) {
	payload := &UpdateGitHubIssueInput{
		Attachments: nil,
	}
	expected, _ := http.NewRequest("POST", "/post", nil)
	expectedForm := url.Values{
		"title": {""},
	}
	expected.Body = io.NopCloser(strings.NewReader(expectedForm.Encode()))
	expected.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	co, err := New(UpdateGitHubIssueInput{})
	assert.NoError(t, err)
	req, err := co.NewRequest("POST", "/post", payload)
	assert.NoError(t, err)
	assert.Equal(t, expected, req)
}

func createTempFile(t *testing.T, content []byte) string {
	t.Helper()
	f, err := os.CreateTemp("", "httpin_test_*.txt")
	assert.NoError(t, err)
	_, err = f.Write(content)
	assert.NoError(t, err)
	return f.Name()
}

func breakMultipartFormBoundary(body *bytes.Buffer) *bytes.Buffer {
	raw := body.Bytes()
	var brokenBody = bytes.NewBuffer(raw[:len(raw)-10])
	brokenBody.Write([]byte("xxx")) // break the boundary
	return brokenBody
}

// assertDecodedFile only checks the File instance that is decoded from the request.
// Don't use it to verify the File instance that is created for upload on the client side.
func assertDecodedFile(t *testing.T, gotFile *File, filename string, content []byte) {
	assert.NotNil(t, gotFile)
	assert.False(t, gotFile.IsUpload())
	assert.Equal(t, gotFile.Filename(), gotFile.Filename())
	assert.Equal(t, filename, gotFile.Filename())
	assert.Equal(t, int64(len(content)), gotFile.Size())
	uploadedContent, err := gotFile.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, content, uploadedContent)
}

func newMultipartFormRequestFromMap(m map[string]any) *http.Request {
	body, writer := newMultipartFormWriterFromMap(m)
	r, _ := http.NewRequest("POST", "/", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	return r
}

func newMultipartFormWriterFromMap(m map[string]any) (body *bytes.Buffer, writer *multipart.Writer) {
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

func removeFileType[T any]() {
	delete(fileTypes, internal.TypeOf[T]())
}
