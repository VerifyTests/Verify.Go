package diff

func defineGuiffy() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp, "-ge2"}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target, "-ge1"}
	}

	return &ToolDefinition{
		Kind:           Guiffy,
		Url:            "https://www.guiffy.com/",
		Cost:           Paid,
		AutoRefresh:    false,
		IsMdi:          false,
		SupportsText:   true,
		RequiresTarget: true,
		BinaryExtensions: []string{
			"bmp", "gif", "jpeg", "jpg", "png", "wbmp",
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/Guiffy/guiffyCL.command",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Guiffy\\guiffy.exe",
			},
		},
		Notes: " * [Command line reference](https://www.guiffy.com/help/GuiffyHelp/GuiffyCmd.html)\n * [Image Diff Tool](https://www.guiffy.com/Image-Diff-Tool.html)\n * `-ge1`: Forbid first file view Editing\n * `-ge2`: Forbid second file view Editing",
	}
}