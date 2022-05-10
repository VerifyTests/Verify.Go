package diff

func defineDiffMerge() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"--nosplash", target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"--nosplash", temp, target}
	}

	return &ToolDefinition{
		Kind:             DiffMerge,
		Url:              "https://www.sourcegear.com/diffmerge/",
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
				"/Applications/DiffMerge.app/Contents/MacOS/DiffMerge",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\SourceGear\\Common\\DiffMerge\\sgdm.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/bin/diffmerge",
			},
		},
	}
}