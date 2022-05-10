package diff

func defineAraxisMerge() *ToolDefinition {
	leftWindowsArgs := func(temp, target string) []string {
		return []string{"/nowait", target, temp}
	}

	rightWindowsArgs := func(temp, target string) []string {
		return []string{"/nowait", temp, target}
	}

	leftMacArgs := func(temp, target string) []string {
		return []string{target, temp}
	}

	rightMacArgs := func(temp, target string) []string {
		return []string{temp, target}
	}

	return &ToolDefinition{
		Kind:           AraxisMerge,
		Url:            "https://www.araxis.com/merge",
		Cost:           Paid,
		AutoRefresh:    true,
		IsMdi:          true,
		SupportsText:   true,
		RequiresTarget: true,
		BinaryExtensions: []string{
			"bmp", "dib", "emf", "gif", "jif", "j2c", "j2k",
			"jp2", "jpc", "jpeg", "jpg", "jpx", "pbm", "pcx",
			"pgm", "png", "ppm", "ras", "tif", "tiff", "tga", "wmf",
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftWindowsArgs,
			TargetRightArguments: rightWindowsArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Araxis\\Araxis Merge\\Compare.exe",
			},
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftMacArgs,
			TargetRightArguments: rightMacArgs,
			ExePaths:             []string{"/Applications/Araxis Merge.app/Contents/Utilities/compare"},
		},
		Notes: "\n" +
			" * [Supported image files](https://www.araxis.com/merge/documentation-windows/comparing-image-files.en)\n" +
			" * [Windows command line usage](https://www.araxis.com/merge/documentation-windows/command-line.en)\n" +
			" * [MacOS command line usage](https://www.araxis.com/merge/documentation-os-x/command-line.en)\n" +
			" * [Installing MacOS command line](https://www.araxis.com/merge/documentation-os-x/installing.en)\n",
		//TODO: add doco about auto refresh
	}

}