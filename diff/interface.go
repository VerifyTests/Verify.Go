package diff

// BuildArguments a function to build arguments to launch a diff tool
type BuildArguments = func(tempFile string, targetFile string) []string

// CIDetected result of detection of Continuous Integration check
type CIDetected bool

const (
	//Detected a Continuous Integration environment
	Detected = CIDetected(true)
	//NotDetected a Continuous Integration environment
	NotDetected = CIDetected(false)
)

// EnvReader reads environment variables
type EnvReader interface {
	LookupEnv(key string) (string, bool)
}

// BuildServerEnvironment environment variable or lookup key to detect CI environment
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

// LaunchResult result of the launch operation of the diff tool
type LaunchResult int

const (
	//NoLaunchResult no result specified
	NoLaunchResult LaunchResult = iota
	//NoEmptyFileForExtension when no empty file can be created for the diff tool
	NoEmptyFileForExtension
	//AlreadyRunningAndSupportsRefresh the diff tool is already running and supports auto refresh
	AlreadyRunningAndSupportsRefresh
	//StartedNewInstance a new instance of the diff tool can be created
	StartedNewInstance
	//TooManyRunningDiffTools too many instance of the diff tool is running
	TooManyRunningDiffTools
	//NoDiffToolFound no diff tool was found on the machine
	NoDiffToolFound
	//Disabled diff tools are disabled
	Disabled
)

//ToolKind specifies the kind of the diff tool detected
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

//PriceModel of the detected diff tool
type PriceModel string

const (
	Paid     PriceModel = "Paid"
	Free     PriceModel = "Free"
	Donation PriceModel = "Free with option to donate"
	Sponsor  PriceModel = "Free with option to sponsor"
)

//AllTools list of all supported diff tools
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

//TryResolveTool a function that possibly finds a diff tool
type TryResolveTool func() (resolved *ResolvedTool, found bool)
