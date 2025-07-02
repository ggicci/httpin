package testutil

import (
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NamedTempFile struct {
	Filename string
	Content  []byte
}

func CreateTempFileV2(t *testing.T) *NamedTempFile {
	t.Helper()
	f, err := os.CreateTemp("", "httpin_test_*.txt")
	assert.NoError(t, err)
	randomContent := RandomString(32)
	_, err = f.Write([]byte(randomContent))
	assert.NoError(t, err)
	f.Close()

	return &NamedTempFile{
		Filename: f.Name(),
		Content:  []byte(randomContent),
	}
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result = make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
