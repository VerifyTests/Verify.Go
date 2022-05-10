package diff

func defineNeovim() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"-d", target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"-d", temp, target}
	}

	return &ToolDefinition{
		Kind:             Neovim,
		Url:              "https://neovim.io/",
		Cost:             Sponsor,
		AutoRefresh:      false,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ChocolateyToolsLocation%\\neovim\\*\\nvim.exe",
			},
		},
		Notes: " * Assumes installed through Chocolatey https://chocolatey.org/packages/neovim/",
	}
}