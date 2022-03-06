package utils

import (
	"fmt"
	"os"
	"strings"
)

var Guard = Guards{}

type Guards struct {
}

// AgainstNullOrEmptySlice guards against the slice being nil or having zero length
func (g *Guards) AgainstNullOrEmptySlice(values []string) {
	if values == nil {
		panic("values provided is nil")
	}

	if len(values) == 0 {
		panic("values cannot be empty")
	}
}

// AgainstEmpty guards against string being empty
func (g *Guards) AgainstEmpty(value string) {
	if len(value) == 0 || value == "" {
		panic("value was expected but was nil or empty")
	}
}

// AgainstBadExtension guards agaist having a "." in the provided file extension.
func (g *Guards) AgainstBadExtension(value string) {
	if strings.HasPrefix(value, ".") {
		panic("Must not start with a period ('.').")
	}
}

// GuardFiles checks if the files exists and the have content
func (g *Guards) GuardFiles(tempFile, targetFile string) {
	g.FileExists(tempFile)
	g.AgainstEmpty(targetFile)
}

// FileExists checks if a file exists at the provided path
func (g *Guards) FileExists(path string) {
	g.AgainstEmpty(path)
	if _, err := os.Stat(path); err != nil {
		panic(fmt.Sprintf("File not found. Path: %s", path))
	}
}
