package diff

import "fmt"

// ResolvedTool contains information about a found diff tool
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

func (r *ResolvedTool) buildCommand(tempFile, targetFile string) string {
	return fmt.Sprintf("\"%s\" %s", r.ExePath, r.getArguments(tempFile, targetFile))
}

func (r *ResolvedTool) commandAndArguments(tempFile, targetFile string) (arguments []string, command string) {
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
