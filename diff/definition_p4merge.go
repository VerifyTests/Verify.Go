package diff

func defineP4MergeImage() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target}
	}

	return &ToolDefinition{
		Kind:           P4MergeImage,
		Url:            "https://www.perforce.com/products/helix-core-apps/merge-diff-tool-p4merge",
		Cost:           Free,
		AutoRefresh:    false,
		IsMdi:          false,
		SupportsText:   false,
		RequiresTarget: true,
		BinaryExtensions: []string{
			"bmp", "gif", "jpg", "jpeg", "png", "pbm",
			"pgm", "ppm", "tif", "tiff", "xbm", "xpm",
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/p4merge.app/Contents/MacOS/p4merge",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Perforce\\p4merge.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/bin/p4merge",
			},
		},
	}
}

func defineP4MergeText() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"-C", "utf8-bom", temp, target, target, target}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"-C", "utf8-bom", target, temp, target, target}
	}

	return &ToolDefinition{
		Kind:             P4MergeText,
		Url:              "https://www.perforce.com/products/helix-core-apps/merge-diff-tool-p4merge",
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
				"/Applications/p4merge.app/Contents/MacOS/p4merge",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Perforce\\p4merge.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/bin/p4merge",
			},
		},
	}
}