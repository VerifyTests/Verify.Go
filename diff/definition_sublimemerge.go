package diff

func defineSublimeMerge() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"mergetool", target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"mergetool", temp, target}
	}

	return &ToolDefinition{
		Kind:             SublimeMerge,
		Url:              "https://www.sublimemerge.com/",
		Cost:             Paid,
		AutoRefresh:      false,
		IsMdi:            true,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/smerge.app/Contents/MacOS/smerge",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Sublime Merge\\smerge.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/bin/smerge",
			},
		},
		Notes: "While SublimeMerge is not MDI, it is treated as MDI since it uses a single shared process to managing multiple windows. As such it is not possible to close a Sublime merge process for a specific diff. [Vote for this feature](https://github.com/sublimehq/sublime_merge/issues/1168)",
	}
}