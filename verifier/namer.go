package verifier

import (
	"fmt"
	"runtime"
	"strings"
)

type namer struct {
	architecture    string
	operatingSystem string
	version         string
	settings        *verifySettings
}

func newNamer(settings *verifySettings) *namer {
	namer := namer{settings: settings}
	return &namer
}

func (n *namer) getUniqueness() string {
	builder := strings.Builder{}
	if n.settings.uniqueForArchitecture {
		builder.WriteString("." + n.getArchitecture())
	}
	if n.settings.uniqueForOperatingSystem {
		builder.WriteString("." + n.getOsPlatform())
	}
	if n.settings.uniqueForRuntime {
		builder.WriteString("." + n.getRuntimeVersion())
	}
	return builder.String()
}

func (n *namer) getArchitecture() string {
	return runtime.GOARCH
}

func (n *namer) getRuntimeVersion() string {
	ver := runtime.Version()
	return strings.ReplaceAll(ver, ".", "_")
}

func (n *namer) getOsPlatform() string {
	if runtime.GOOS == "windows" {
		return "windows"
	}

	if runtime.GOOS == "unix" {
		return "linux"
	}

	if runtime.GOOS == "linux" {
		return "linux"
	}

	if runtime.GOOS == "darwin" {
		return "macos"
	}

	panic(fmt.Sprintf("Unknown OS: %s", runtime.GOOS))
}
