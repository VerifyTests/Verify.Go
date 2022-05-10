package diff

func defineKDiff3() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp, "--cs", "CreateBakFiles=0"}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target, "--cs", "CreateBakFiles=0"}
	}

	return &ToolDefinition{
		Kind:             KDiff3,
		Url:              "https://github.com/KDE/kdiff3",
		Cost:             Free,
		AutoRefresh:      false,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/kdiff3.app/Contents/MacOS/kdiff3",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\KDiff3\\kdiff3.exe",
				"%ProgramFiles%\\KDiff3\\bin\\kdiff3.exe",
			},
		},
		Notes: " * `--cs CreateBakFiles=0` to not save a `.orig` file when merging",
	}
}