package verifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileExtension(t *testing.T) {
	assert.Equal(t, "txt", file.getFileExtension("textfile.txt"))
	assert.Equal(t, "json", file.getFileExtension("textfile.json"))
	assert.Equal(t, "txt", file.getFileExtension("txt"))
}

func TestFileReading(t *testing.T) {
	content := file.readText("../_testdata/verifier_test.TestNilTargets.verified.txt")
	assert.NotEmpty(t, content)
}
