package diff

func defineTortoiseGitMerge() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target}
	}

	return &ToolDefinition{
		Kind:             TortoiseGitMerge,
		Url:              "https://tortoisegit.org/docs/tortoisegitmerge/",
		Cost:             Free,
		AutoRefresh:      false,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\TortoiseGit\\bin\\TortoiseGitMerge.exe",
			},
		},
	}
}