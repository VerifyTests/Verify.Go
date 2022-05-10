package diff

func defineGoland() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"diff", target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"diff", temp, target}
	}

	return &ToolDefinition{
		Kind:             GoLand,
		Url:              "https://www.jetbrains.com/go/",
		Cost:             Paid,
		AutoRefresh:      false,
		IsMdi:            true,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%HOME%/Library/Application Support/JetBrains/Toolbox/apps/GoLand/*/*/GoLand EAP.app/Contents/MacOS/goland",
				"%HOME%/Library/Application Support/JetBrains/Toolbox/apps/GoLand/*/*/GoLand.app/Contents/MacOS/goland",
				"/Applications/GoLand EAP.app/Contents/MacOS/goland",
				"/Applications/GoLand.app/Contents/MacOS/goland",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%LOCALAPPDATA%\\JetBrains\\Installations\\GoLand*\\bin\\goland64.exe",
				"%ProgramFiles%\\JetBrains\\JetBrains GoLand *\\bin\\goland64.exe",
				"%JetBrains GoLand%\\goland64.exe",
				"%LOCALAPPDATA%\\JetBrains\\Toolbox\\apps\\GoLand\\*\\*\\bin\\goland64.exe",
				"%UserProfile%\\scoop\\apps\\goland\\current\\IDE\\bin\\goland64.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%HOME%/.local/share/JetBrains/Toolbox/apps/Rider/*/*/bin/rider.sh",
				"/opt/jetbrains/rider/bin/rider.sh",
				"/usr/share/rider/bin/rider.sh",
			},
		},
		Notes: " * https://www.jetbrains.com/help/rider/Command_Line_Differences_Viewer.html",
	}
}