package diff

func defineVsCode() *ToolDefinition {
	leftArg := func(temp, target string) []string {
		return []string{"--diff", target, temp}
	}

	rightArg := func(temp, target string) []string {
		return []string{"--diff", temp, target}
	}

	return &ToolDefinition{
		Kind:             VisualStudioCode,
		Url:              "https://code.visualstudio.com",
		Cost:             Free,
		AutoRefresh:      true,
		IsMdi:            true,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: make([]string, 0),
		Windows: OsSettings{
			TargetLeftArguments:  leftArg,
			TargetRightArguments: rightArg,
			ExePaths: []string{
				"%LocalAppData%\\Programs\\Microsoft VS Code\\code.exe",
				"%ProgramFiles%\\Microsoft VS Code\\bin\\code.exe",
				"%ProgramFiles%\\Microsoft VS Code\\code.exe",
				"%UserProfile%\\scoop\\apps\\vscode\\current\\bin\\code.cmd",
				"%UserProfile%\\scoop\\apps\\vscode\\current\\code.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArg,
			TargetRightArguments: rightArg,
			ExePaths: []string{
				"/usr/local/bin/code",
				"/usr/bin/code",
			},
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftArg,
			TargetRightArguments: rightArg,
			ExePaths: []string{
				"/Applications/Visual Studio Code.app/Contents/MacOS/Electron",
				"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
			},
		},
		Notes: "\n* [Command line reference](https://code.visualstudio.com/docs/editor/command-line)",
	}
}