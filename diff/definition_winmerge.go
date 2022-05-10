package diff

import "github.com/VerifyTests/Verify.Go/utils"

func defineWinMerge() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		tempTitle := utils.File.GetFileName(temp)
		targetTitle := utils.File.GetFileName(target)
		return []string{"/u", "/wl", "/e", target, temp, "/dl", targetTitle, "/dr", tempTitle}
	}

	rightArgs := func(temp, target string) []string {
		tempTitle := utils.File.GetFileName(temp)
		targetTitle := utils.File.GetFileName(target)
		return []string{"/u", "/wl", "/e", temp, target, "/dl", tempTitle, "/dr", targetTitle}
	}

	return &ToolDefinition{
		Kind:           WinMerge,
		Url:            "https://winmerge.org/",
		Cost:           Donation,
		AutoRefresh:    true,
		IsMdi:          false,
		SupportsText:   true,
		RequiresTarget: true,
		BinaryExtensions: []string{
			"bmp", "cut", "dds", "exr", "g3", "gif", "hdr", "ico",
			"iff", "lbm", "j2k", "j2c", "jng", "jp2", "jpg", "jif",
			"jpeg", "jpe", "jxr", "wdp", "hdp", "koa", "mng", "pcd",
			"pcx", "pfm", "pct", "pict", "pic", "png", "pbm", "pgm",
			"ppm", "psd", "ras", "sgi", "rgb", "rgba", "bw", "tga",
			"targa", "tif", "tiff", "wap", "wbmp", "wbm", "webp", "xbm", "xpm",
		},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\WinMerge\\WinMergeU.exe",
				"%LocalAppData%\\Programs\\WinMerge\\WinMergeU.exe",
			},
		},
		Notes: `
 * [Command line reference](https://manual.winmerge.org/en/Command_line.html).
 * '/u'' Prevents WinMerge from adding paths to the Most Recently Used (MRU) list.
 * '/wl' Opens the left side as read-only.
 * '/dl' and '/dr' Specifies file descriptions in the title bar.
 * '/e' Enables close with a single Esc key press.`,
	}
}