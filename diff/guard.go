package diff

import (
	"fmt"
	"os"
)

var guard = guardChecker{}

type guardChecker struct {
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

func (g *guardChecker) AgainstEmpty(value string) {
	if len(value) == 0 || value == "" {
		panic("value was expected but was nil or empty")
	}
}
