package diff

type BuildArguments = func(tempFile string, targetFile string) []string
type CIDetected bool

const (
	CI_Detected    = CIDetected(true)
	CI_NotDetected = CIDetected(false)
)

type EnvReader interface {
	LookupEnv(key string) (string, bool)
}

type BuildServerEnvironment string

const (
	Appveyor    BuildServerEnvironment = "CI"
	Jenkins     BuildServerEnvironment = "JENKINS_URL"
	GitHub      BuildServerEnvironment = "GITHUB_ACTION"
	AzureDevOps BuildServerEnvironment = "AGENT_ID"
	TeamCity    BuildServerEnvironment = "TEAMCITY_VERSION"
	MyGet       BuildServerEnvironment = "BuildRunner"
	GitLab      BuildServerEnvironment = "GITLAB_CI"
)

type LaunchResult int

const (
	NoLaunchResult LaunchResult = iota
	NoEmptyFileForExtension
	AlreadyRunningAndSupportsRefresh
	StartedNewInstance
	TooManyRunningDiffTools
	NoDiffToolFound
	Disabled
)

type ToolKind string

const (
	None             ToolKind = ""
	BeyondCompare    ToolKind = "BeyondCompare"
	P4MergeText      ToolKind = "P4MergeText"
	P4MergeImage     ToolKind = "P4MergeImage"
	AraxisMerge      ToolKind = "AraxisMerge"
	Meld             ToolKind = "Meld"
	SublimeMerge     ToolKind = "SublimeMerge"
	Kaleidoscope     ToolKind = "Kaleidoscope"
	CodeCompare      ToolKind = "CodeCompare"
	DeltaWalker      ToolKind = "DeltaWalker"
	WinMerge         ToolKind = "WinMerge"
	DiffMerge        ToolKind = "DiffMerge"
	TortoiseMerge    ToolKind = "TortoiseMerge"
	TortoiseGitMerge ToolKind = "TortoiseGitMerge"
	TortoiseIDiff    ToolKind = "TortoiseIDiff"
	KDiff3           ToolKind = "KDiff3"
	TkDiff           ToolKind = "TkDiff"
	Guiffy           ToolKind = "Guiffy"
	ExamDiff         ToolKind = "ExamDiff"
	Diffinity        ToolKind = "Diffinity"
	VisualStudioCode ToolKind = "VisualStudioCode"
	GoLand           ToolKind = "GoLand"
	Vim              ToolKind = "Vim"
	Neovim           ToolKind = "Neovim"
)

type PriceModel string

const (
	Paid     PriceModel = "Paid"
	Free     PriceModel = "Free"
	Donation PriceModel = "Free with option to donate"
	Sponsor  PriceModel = "Free with option to sponsor"
)

var AllTools = []ToolKind{
	BeyondCompare,
	P4MergeText,
	P4MergeImage,
	AraxisMerge,
	Meld,
	SublimeMerge,
	Kaleidoscope,
	CodeCompare,
	DeltaWalker,
	WinMerge,
	DiffMerge,
	TortoiseMerge,
	TortoiseGitMerge,
	TortoiseIDiff,
	KDiff3,
	TkDiff,
	Guiffy,
	ExamDiff,
	Diffinity,
	VisualStudioCode,
	GoLand,
	Vim,
	Neovim,
}

type TryResolveTool func() (resolved *ResolvedTool, found bool)
