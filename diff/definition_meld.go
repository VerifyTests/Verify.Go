package diff

func defineMeld() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target}
	}

	return &ToolDefinition{
		Kind:             Meld,
		Url:              "https://meldmerge.org/",
		Cost:             Free,
		AutoRefresh:      false,
		IsMdi:            true,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/meld.app/Contents/MacOS/meld",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%LOCALAPPDATA%\\Programs\\Meld\\meld.exe",
				"%ProgramFiles%\\Meld\\meld.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/bin/meld",
			},
		},
		Notes: "While Meld is not MDI, it is treated as MDI since it uses a single shared process to managing multiple windows. As such it is not possible to close a Meld merge process for a specific diff. [Vote for this feature](https://gitlab.gnome.org/GNOME/meld/-/issues/584)",
	}
}