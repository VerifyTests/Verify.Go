package diff

func defineDeltaWalker() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"-mi", target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"-mi", temp, target}
	}

	return &ToolDefinition{
		Kind:           DeltaWalker,
		Url:            "https://www.deltawalker.com/",
		Cost:           Paid,
		AutoRefresh:    false,
		IsMdi:          false,
		SupportsText:   true,
		RequiresTarget: false,
		BinaryExtensions: []string{
			"jpg", "jp2", "j2k", "png", "gif", "psd", "tif", "bmp",
			"pct", "pict", "pic", "ico", "ppm", "pgm", "pbm", "pnm",
			"zip", "jar", "ear", "tar", "tgz", "tbz2", "gz", "bz2", "doc",
			"docx", "xls", "xlsx", "ppt", "pdf", "rtf", "html", "htm",
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/DeltaWalker.app/Contents/MacOS/DeltaWalker",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"C:\\Program Files\\Deltopia\\DeltaWalker\\DeltaWalker.exe",
			},
		},
		Notes: " * [Command line usage](https://www.deltawalker.com/integrate/command-line)",
	}
}