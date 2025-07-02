package testutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	assert.Len(t, RandomString(32), 32)
}

func TestCreateTempFileV2(t *testing.T) {
	tempFile := CreateTempFileV2(t)
	assert.NotEmpty(t, tempFile.Filename)
	assert.Len(t, tempFile.Content, 32)
	assert.Contains(t, tempFile.Filename, "httpin_test_")
	assert.FileExists(t, tempFile.Filename)

	// Clean up the temporary file after the test
	err := os.Remove(tempFile.Filename)
	assert.NoError(t, err)
}
