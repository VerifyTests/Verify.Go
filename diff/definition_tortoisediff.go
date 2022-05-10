package diff

func defineTortoiseDiff() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"/left:" + target, "/right:" + temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"/left:" + temp, "/right:" + target}
	}

	return &ToolDefinition{
		Kind:           TortoiseIDiff,
		Url:            "https://tortoisesvn.net/TortoiseIDiff.html",
		Cost:           Free,
		AutoRefresh:    false,
		IsMdi:          false,
		SupportsText:   false,
		RequiresTarget: true,
		BinaryExtensions: []string{
			"bmp", "gif", "jpg", "jpeg",
			"png", "ico", "tif", "tiff",
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\TortoiseSVN\\bin\\TortoiseIDiff.exe",
			},
		},
	}
}