package diff

func defineVim() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		return []string{"-d", target, temp, "-c", "setl autoread | setl nobackup | set noswapfile"}
	}

	rightArgs := func(temp, target string) []string {
		return []string{"-d", temp, target, "-c", "setl autoread | setl nobackup | set noswapfile"}
	}

	return &ToolDefinition{
		Kind:             Vim,
		Url:              "https://www.vim.org/",
		Cost:             Donation,
		AutoRefresh:      true,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\Vim\\*\\vim.exe",
			},
		},
		Osx: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"/Applications/MacVim.app/Contents/bin/mvim",
			},
		},
		Notes: `
* [Options](http://vimdoc.sourceforge.net/htmldoc/options.html)
* [Vim help files](https://vimhelp.org/)
* [autoread](http://vimdoc.sourceforge.net/htmldoc/options.html#'autoread')
* [nobackup](http://vimdoc.sourceforge.net/htmldoc/options.html#'backup')
* [noswapfile](http://vimdoc.sourceforge.net/htmldoc/options.html#'swapfile')`,
	}
}