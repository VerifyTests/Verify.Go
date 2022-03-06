package utils

import (
	"fmt"
	"os"
	"strings"
)

var Guard = guardChecker{}

type guardChecker struct {
}

func (g *guardChecker) AgainstNullOrEmptySlice(values []string) {
	if values == nil {
		panic("values provided is nil")
	}

	if len(values) == 0 {
		panic("values cannot be empty")
	}
}

func (g *guardChecker) AgainstEmpty(value string) {
	if len(value) == 0 || value == "" {
		panic("value was expected but was nil or empty")
	}
}

func (g *guardChecker) AgainstBadExtension(value string) {
	if strings.HasPrefix(value, ".") {
		panic("Must not start with a period ('.').")
	}
}

func (g *guardChecker) GuardFiles(tempFile, targetFile string) {
	g.FileExists(tempFile)
	g.AgainstEmpty(targetFile)
}

func (g *guardChecker) FileExists(path string) {
	g.AgainstEmpty(path)
	if _, err := os.Stat(path); err != nil {
		panic(fmt.Sprintf("File not found. Path: %s", path))
	}
}
