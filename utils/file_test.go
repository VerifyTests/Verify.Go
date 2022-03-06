package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileExtension(t *testing.T) {
	assert.Equal(t, "txt", File.GetFileExtension("textfile.txt"))
	assert.Equal(t, "json", File.GetFileExtension("textfile.json"))
	assert.Equal(t, "txt", File.GetFileExtension("txt"))
}

func TestFileReading(t *testing.T) {
	content := File.ReadText("../_testdata/verifier_test.TestNilTargets.verified.txt")
	assert.NotEmpty(t, content)
}
