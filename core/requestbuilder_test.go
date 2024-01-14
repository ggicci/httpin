package core

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// FIX: https://github.com/ggicci/httpin/issues/88
// Impossible to make streaming
func TestIssue88_RequestBuilderFileUploadStreaming(t *testing.T) {
	rb := NewRequestBuilder(context.Background())

	var contentReader = strings.NewReader("hello world")
	rb.SetAttachment("file", []FileMarshaler{
		UploadStream(io.NopCloser(contentReader)),
	})
	req, _ := http.NewRequest("GET", "/", nil)
	rb.Populate(req)

	err := req.ParseMultipartForm(32 << 20)
	assert.NoError(t, err)

	file, fh, err := req.FormFile("file")
	assert.NoError(t, err)
	assert.Equal(t, "file_0", fh.Filename)
	content, err := io.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(content))
}

func TestIssue88_CancelStreaming(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	rb := NewRequestBuilder(ctx)
	var contentReader = strings.NewReader("hello world")
	rb.SetAttachment("file", []FileMarshaler{
		UploadStream(io.NopCloser(contentReader)),
	})
	req, _ := http.NewRequest("GET", "/", nil)
	rb.Populate(req)

	cancel()
	time.Sleep(time.Millisecond * 100)

	err := req.ParseMultipartForm(32 << 20)
	assert.ErrorContains(t, err, "context canceled")
}
