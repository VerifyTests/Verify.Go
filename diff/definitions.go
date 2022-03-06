package diff

import "github.com/VerifyTests/Verify.Go/utils"

type ToolDefinition struct {
	Kind             ToolKind
	Url              string
	AutoRefresh      bool
	IsMdi            bool
	Windows          OsSettings
	Linux            OsSettings
	Osx              OsSettings
	BinaryExtensions []string
	Cost             PriceModel
	Notes            string
	SupportsText     bool
	RequiresTarget   bool
}

func newToolDefinition(tool ToolKind, url string, autoRefresh bool,
	isMdi bool, windows OsSettings, linux OsSettings, osx OsSettings,
	binaryExtensions []string, cost string, notes string, supportsText bool, requiresTarget bool) ToolDefinition {
	return ToolDefinition{
		Kind:             tool,
		Url:              url,
		AutoRefresh:      autoRefresh,
		IsMdi:            isMdi,
		Windows:          windows,
		Linux:            linux,
		Osx:              osx,
		BinaryExtensions: binaryExtensions,
		Cost:             PriceModel(cost),
		Notes:            notes,
		SupportsText:     supportsText,
		RequiresTarget:   requiresTarget}
}

var AllDefinedTools = []*ToolDefinition{
	defineBeyondCompare(),
	defineVsCode(),
	defineAraxisMerge(),
	defineCodeCompare(),
	defineDeltaWalker(),
	defineDiffinity(),
	defineDiffMerge(),
	defineExamDiff(),
	defineGuiffy(),
	defineKaleidoscope(),
	defineKDiff3(),
	defineMeld(),
	defineNeovim(),
	defineP4MergeImage(),
	defineP4MergeText(),
	defineRider(),
	defineSublimeMerge(),
}

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

func defineDiffinity() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target}
	}

	return &ToolDefinition{
		Kind:             Diffinity,
		Url:              "https://truehumandesign.se/s_diffinity.php",
		Cost:             Donation,
		AutoRefresh:      false,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Diffinity\\Diffinity.exe",
				"%UserProfile%\\scoop\\apps\\diffinity\\current\\Diffinity.exe",
			},
		},
		Notes: "* Disable single instance:\n   \\ Preferences \\ Tabs \\ uncheck `Use single instance and open new diffs in tabs`.",
	}
}

func defineDiffMerge() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"--nosplash", target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"--nosplash", temp, target}
	}

	return &ToolDefinition{
		Kind:             DiffMerge,
		Url:              "https://www.sourcegear.com/diffmerge/",
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
				"/Applications/DiffMerge.app/Contents/MacOS/DiffMerge",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\SourceGear\\Common\\DiffMerge\\sgdm.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/bin/diffmerge",
			},
		},
	}
}

func defineExamDiff() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		tempTitle := utils.File.GetFileName(temp)
		targetTitle := utils.File.GetFileName(target)
		return []string{target, temp, "/nh", "/diffonly", "/dn1:" + targetTitle, "/dn2:" + tempTitle}
	}

	rightArgs := func(temp, target string) []string {
		tempTitle := utils.File.GetFileName(temp)
		targetTitle := utils.File.GetFileName(target)
		return []string{temp, target, "/nh", "/diffonly", "/dn1:" + tempTitle, "/dn2:" + targetTitle}
	}

	return &ToolDefinition{
		Kind:             ExamDiff,
		Url:              "https://www.prestosoft.com/edp_examdiffpro.asp",
		Cost:             Paid,
		AutoRefresh:      true,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\ExamDiff Pro\\ExamDiff.exe",
			},
		},
		Notes: " * [Command line reference](https://www.prestosoft.com/ps.asp?page=htmlhelp/edp/command_line_options)\n * `/nh`: do not add files or directories to comparison history\n * `/diffonly`: diff-only merge mode: hide the merge pane",
	}
}

func defineGuiffy() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp, "-ge2"}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target, "-ge1"}
	}

	return &ToolDefinition{
		Kind:           Guiffy,
		Url:            "https://www.guiffy.com/",
		Cost:           Paid,
		AutoRefresh:    false,
		IsMdi:          false,
		SupportsText:   true,
		RequiresTarget: true,
		BinaryExtensions: []string{
			"bmp", "gif", "jpeg", "jpg", "png", "wbmp",
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/Guiffy/guiffyCL.command",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Guiffy\\guiffy.exe",
			},
		},
		Notes: " * [Command line reference](https://www.guiffy.com/help/GuiffyHelp/GuiffyCmd.html)\n * [Image Diff Tool](https://www.guiffy.com/Image-Diff-Tool.html)\n * `-ge1`: Forbid first file view Editing\n * `-ge2`: Forbid second file view Editing",
	}
}

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

func defineKDiff3() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp, "--cs", "CreateBakFiles=0"}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target, "--cs", "CreateBakFiles=0"}
	}

	return &ToolDefinition{
		Kind:             KDiff3,
		Url:              "https://github.com/KDE/kdiff3",
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
				"/Applications/kdiff3.app/Contents/MacOS/kdiff3",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\KDiff3\\kdiff3.exe",
				"%ProgramFiles%\\KDiff3\\bin\\kdiff3.exe",
			},
		},
		Notes: " * `--cs CreateBakFiles=0` to not save a `.orig` file when merging",
	}
}

func defineMeld() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{temp, target}
	}

	return &ToolDefinition{
		Kind:             Meld,
		Url:              "https://meldmerge.org/",
		Cost:             Free,
		AutoRefresh:      false,
		IsMdi:            true,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/meld.app/Contents/MacOS/meld",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%LOCALAPPDATA%\\Programs\\Meld\\meld.exe",
				"%ProgramFiles%\\Meld\\meld.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/bin/meld",
			},
		},
		Notes: "While Meld is not MDI, it is treated as MDI since it uses a single shared process to managing multiple windows. As such it is not possible to close a Meld merge process for a specific diff. [Vote for this feature](https://gitlab.gnome.org/GNOME/meld/-/issues/584)",
	}
}

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

func defineRider() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"diff", target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"diff", temp, target}
	}

	return &ToolDefinition{
		Kind:             Rider,
		Url:              "https://www.jetbrains.com/rider/",
		Cost:             Paid,
		AutoRefresh:      false,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%HOME%/Library/Application Support/JetBrains/Toolbox/apps/Rider/*/*/Rider EAP.app/Contents/MacOS/rider",
				"%HOME%/Library/Application Support/JetBrains/Toolbox/apps/Rider/*/*/Rider.app/Contents/MacOS/rider",
				"/Applications/Rider EAP.app/Contents/MacOS/rider",
				"/Applications/Rider.app/Contents/MacOS/rider",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%LOCALAPPDATA%\\JetBrains\\Installations\\Rider*\\bin\\rider64.exe",
				"%ProgramFiles%\\JetBrains\\JetBrains Rider *\\bin\\rider64.exe",
				"%JetBrains Rider%\\rider64.exe",
				"%LOCALAPPDATA%\\JetBrains\\Toolbox\\apps\\Rider\\*\\*\\bin\\rider64.exe",
				"%UserProfile%\\scoop\\apps\\rider\\current\\IDE\\bin\\rider64.exe",
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

func defineSublimeMerge() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"mergetool", target, temp}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"mergetool", temp, target}
	}

	return &ToolDefinition{
		Kind:             SublimeMerge,
		Url:              "https://www.sublimemerge.com/",
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
				"/Applications/smerge.app/Contents/MacOS/smerge",
			},
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Sublime Merge\\smerge.exe",
			},
		},
		Linux: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/usr/bin/smerge",
			},
		},
		Notes: "While SublimeMerge is not MDI, it is treated as MDI since it uses a single shared process to managing multiple windows. As such it is not possible to close a Sublime merge process for a specific diff. [Vote for this feature](https://github.com/sublimehq/sublime_merge/issues/1168)",
	}
}
