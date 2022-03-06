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
	//Appveyor build environment
	Appveyor BuildServerEnvironment = "CI"
	//Jenkins build environment
	Jenkins BuildServerEnvironment = "JENKINS_URL"
	//GitHub actions build environment
	GitHub BuildServerEnvironment = "GITHUB_ACTION"
	//AzureDevOps build environment
	AzureDevOps BuildServerEnvironment = "AGENT_ID"
	//TeamCity build environment
	TeamCity BuildServerEnvironment = "TEAMCITY_VERSION"
	//MyGet build environment
	MyGet BuildServerEnvironment = "BuildRunner"
	//GitLab build environment
	GitLab BuildServerEnvironment = "GITLAB_CI"
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
	//None no tool
	None ToolKind = ""
	//BeyondCompare diff tool
	BeyondCompare ToolKind = "BeyondCompare"
	//P4MergeText diff tool
	P4MergeText ToolKind = "P4MergeText"
	//P4MergeImage diff tool
	P4MergeImage ToolKind = "P4MergeImage"
	//AraxisMerge diff tool
	AraxisMerge ToolKind = "AraxisMerge"
	//Meld diff tool
	Meld ToolKind = "Meld"
	//SublimeMerge diff tool
	SublimeMerge ToolKind = "SublimeMerge"
	//Kaleidoscope diff tool
	Kaleidoscope ToolKind = "Kaleidoscope"
	//CodeCompare diff tool
	CodeCompare ToolKind = "CodeCompare"
	//DeltaWalker diff tool
	DeltaWalker ToolKind = "DeltaWalker"
	//WinMerge diff tool
	WinMerge ToolKind = "WinMerge"
	//DiffMerge diff tool
	DiffMerge ToolKind = "DiffMerge"
	//TortoiseMerge diff tool
	TortoiseMerge ToolKind = "TortoiseMerge"
	//TortoiseGitMerge diff tool
	TortoiseGitMerge ToolKind = "TortoiseGitMerge"
	//TortoiseIDiff diff tool
	TortoiseIDiff ToolKind = "TortoiseIDiff"
	//KDiff3 diff tool
	KDiff3 ToolKind = "KDiff3"
	//TkDiff diff tool
	TkDiff ToolKind = "TkDiff"
	//Guiffy diff tool
	Guiffy ToolKind = "Guiffy"
	//ExamDiff diff tool
	ExamDiff ToolKind = "ExamDiff"
	//Diffinity diff tool
	Diffinity ToolKind = "Diffinity"
	//VisualStudioCode diff tool
	VisualStudioCode ToolKind = "VisualStudioCode"
	//GoLand diff tool
	GoLand ToolKind = "GoLand"
	//Vim diff tool
	Vim ToolKind = "Vim"
	//Neovim diff tool
	Neovim ToolKind = "Neovim"
)

//PriceModel of the detected diff tool
type PriceModel string

const (
	//Paid diff tool
	Paid PriceModel = "Paid"
	//Free diff tool
	Free PriceModel = "Free"
	//Donation free with donation options
	Donation PriceModel = "Free with option to donate"
	//Sponsor free with ability to sponsor
	Sponsor PriceModel = "Free with option to sponsor"
)

var allTools = []ToolKind{
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
