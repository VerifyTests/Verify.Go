package diff

import "fmt"

type ResolvedTool struct {
	Name             string
	Kind             ToolKind
	ExePath          string
	RightArguments   BuildArguments
	LeftArguments    BuildArguments
	IsMdi            bool
	AutoRefresh      bool
	BinaryExtensions []string
	RequiresTarget   bool
	SupportsText     bool
}

func newResolvedTool(
	name string,
	tool ToolKind,
	exePath string,
	targetRightArguments BuildArguments,
	targetLeftArguments BuildArguments,
	isMdi bool,
	autoRefresh bool,
	binaryExtensions []string,
	requiresTarget bool,
	supportsText bool) ResolvedTool {

	return ResolvedTool{
		Name:             name,
		Kind:             tool,
		ExePath:          exePath,
		RightArguments:   targetRightArguments,
		LeftArguments:    targetLeftArguments,
		IsMdi:            isMdi,
		AutoRefresh:      autoRefresh,
		BinaryExtensions: binaryExtensions,
		RequiresTarget:   requiresTarget,
		SupportsText:     supportsText,
	}
}

func (r *ResolvedTool) BuildCommand(tempFile, targetFile string) string {
	return fmt.Sprintf("\"%s\" %s", r.ExePath, r.getArguments(tempFile, targetFile))
}

func (r *ResolvedTool) CommandAndArguments(tempFile, targetFile string) (arguments []string, command string) {
	arguments = r.getArguments(tempFile, targetFile)
	command = r.ExePath
	return
}

func (r *ResolvedTool) getArguments(tempFile, targetFile string) []string {
	if position.TargetOnLeft {
		return r.LeftArguments(tempFile, targetFile)
	}
	return r.RightArguments(tempFile, targetFile)
}
