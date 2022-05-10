package diff

func defineKaleidoscope() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target}
	}

	return &ToolDefinition{
		Kind:           Kaleidoscope,
		Url:            "https://www.kaleidoscopeapp.com/",
		Cost:           Paid,
		AutoRefresh:    false,
		IsMdi:          false,
		SupportsText:   true,
		RequiresTarget: true,
		BinaryExtensions: []string{
			"bmp", "gif", "ico", "jpg", "jpeg", "png", "tiff", "tif",
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/local/bin/ksdiff",
			},
		},
	}
}