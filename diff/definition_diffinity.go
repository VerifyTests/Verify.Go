package diff

func defineDiffinity() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target}
	}

	return &ToolDefinition{
		Kind:             Diffinity,
		Url:              "https://truehumandesign.se/s_diffinity.php",
		Cost:             Donation,
		AutoRefresh:      false,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Diffinity\\Diffinity.exe",
				"%UserProfile%\\scoop\\apps\\diffinity\\current\\Diffinity.exe",
			},
		},
		Notes: "* Disable single instance:\n   \\ Preferences \\ Tabs \\ uncheck `Use single instance and open new diffs in tabs`.",
	}
}