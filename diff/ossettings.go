package diff

import (
	"runtime"
	"strings"
)

type OsSettings struct {
	TargetLeftArguments  BuildArguments
	TargetRightArguments BuildArguments
	ExePaths             []string
}

func NewOsSettings(
	targetLeftArguments BuildArguments,
	targetRightArguments BuildArguments,
	exePaths []string) *OsSettings {

	return &OsSettings{TargetLeftArguments: targetLeftArguments, TargetRightArguments: targetRightArguments, ExePaths: exePaths}
}

func resolveFromOsSettings(windows, linux, osx *OsSettings) (path string, leftArguments, rightArguments BuildArguments, found bool) {

	if windows != nil && runtime.GOOS == "windows" {
		paths := expandProgramFiles(windows.ExePaths)
		path, found := finder.TryFindExe(paths)
		if found {
			targetLeftArguments := windows.TargetLeftArguments
			targetRightArguments := windows.TargetRightArguments
			return path, targetLeftArguments, targetRightArguments, true
		}
	}

	if linux != nil && runtime.GOOS == "linux" {
		path, found := finder.TryFindExe(linux.ExePaths)
		if found {
			return path, linux.TargetLeftArguments, linux.TargetRightArguments, true
		}
	}

	if osx != nil && runtime.GOOS == "darwin" {
		path, found := finder.TryFindExe(osx.ExePaths)
		if found {
			return path, osx.TargetLeftArguments, osx.TargetRightArguments, true
		}
	}

	return "", nil, nil, false
}

func expandProgramFiles(paths []string) []string {
	result := make([]string, 0)
	for _, windowsPath := range paths {
		result = append(result, windowsPath)
		if strings.ContainsAny(windowsPath, "%ProgramFiles%") {
			result = append(result, strings.Replace(windowsPath, "%ProgramFiles%", "%ProgramW6432%", 1))
			result = append(result, strings.Replace(windowsPath, "%ProgramFiles%", "%ProgramW6432%", 1))
		}
	}
	return result
}
