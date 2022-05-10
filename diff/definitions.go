package diff

//ToolDefinition specifies attributes of a diff tool
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

//AllDefinedTools definition of all supported diff tools
var AllDefinedTools = []*ToolDefinition{
	defineBeyondCompare(),
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
	defineGoland(),
	defineSublimeMerge(),
	defineTkDiff(),
	defineTortoiseGitMerge(),
	defineTortoiseDiff(),
	defineTortoiseMerge(),
	defineVim(),
	defineVsCode(),
	defineWinMerge(),
}