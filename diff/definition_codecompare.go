package diff

func defineCodeCompare() *ToolDefinition {
	return &ToolDefinition{
		Kind:             CodeCompare,
		Url:              "https://www.devart.com/codecompare/",
		Cost:             Paid,
		AutoRefresh:      false,
		IsMdi:            true,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Windows: OsSettings{
			TargetLeftArguments: func(temp, target string) []string {
				return []string{target, temp}
			},
			TargetRightArguments: func(temp, target string) []string {
				return []string{temp, target}
			},
			ExePaths: []string{
				"%ProgramFiles%\\Devart\\Code Compare\\CodeCompare.exe",
			},
		},
		Notes: "* [Command line reference](https://www.devart.com/codecompare/docs/index.html?comparing_via_command_line.htm)",
	}
}