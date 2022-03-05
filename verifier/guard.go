package verifier

import (
	"strings"
)

var guard = guardChecker{}

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

func (g *guardChecker) AgainstNullOrEmpty(value string) {
	if len(value) == 0 || value == "" {
		panic("value was expected but was nil or empty")
	}
}

func (g *guardChecker) AgainstBadExtension(value string) {
	if strings.HasPrefix(value, ".") {
		panic("Must not start with a period ('.').")
	}
}
