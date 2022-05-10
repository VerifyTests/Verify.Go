package diff

func defineBeyondCompare() *ToolDefinition {
	leftWindowsArgs := func(temp, target string) []string {
		return []string{"/solo", "/rightreadonly", target, temp}
	}

	rightWindowsArgs := func(temp, target string) []string {
		return []string{"/solo", "/leftreadonly", temp, target}
	}

	leftLinuxArgs := func(temp, target string) []string {
		return []string{"-solo", "-rightreadonly", target, temp}
	}

	rightLinuxArgs := func(temp, target string) []string {
		return []string{"-solo", "-leftreadonly", temp, target}
	}

	return &ToolDefinition{
		Kind:           BeyondCompare,
		Url:            "https://www.scootersoftware.com",
		Cost:           Paid,
		AutoRefresh:    true,
		IsMdi:          false,
		SupportsText:   true,
		RequiresTarget: true,
		BinaryExtensions: []string{
			"mp3", "xls", "xlsm", "xlsx", "doc", "docm", "docx", "dot", "dotm", "dotx",
			"pdf", "bmp", "gif", "ico", "jpg", "jpeg", "png", "tif", "tiff", "rtf",
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftWindowsArgs,
			TargetRightArguments: rightWindowsArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Beyond Compare *\\BCompare.exe",
				"%UserProfile%\\scoop\\apps\\beyondcompare\\current\\BCompare.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftLinuxArgs,
			TargetRightArguments: rightLinuxArgs,
			ExePaths:             []string{"/usr/lib/beyondcompare/bcomp"},
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftLinuxArgs,
			TargetRightArguments: rightLinuxArgs,
			ExePaths:             []string{"/Applications/Beyond Compare.app/Contents/MacOS/bcomp"},
		},
		Notes: "* [Command line reference](https://www.scootersoftware.com/v4help/index.html?command_line_reference.html)",
	}
}