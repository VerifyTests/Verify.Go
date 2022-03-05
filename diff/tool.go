package diff

import (
	"fmt"
	"os"
	"strings"
)

type tools struct {
	resolved        []*ResolvedTool
	pathLookup      map[string]*ResolvedTool
	extensionLookup map[string]*ResolvedTool
}

func newTools() *tools {
	t := &tools{}
	t.reset()
	return t
}

func (t *tools) tryFind(kind ToolKind) (tool *ResolvedTool, found bool) {
	for _, rt := range t.resolved {
		if rt.Kind == kind {
			return rt, true
		}
	}
	return nil, false
}

func (t *tools) tryFindForExtension(extension string) (tool *ResolvedTool, found bool) {
	extension = file.getFileExtension(extension)
	if file.isText(extension) {
		for _, tool := range t.resolved {
			if tool.SupportsText {
				return tool, true
			}
		}

		tool, found = t.extensionLookup[extension]
		return
	}
	return nil, false
}

func (t *tools) reset() {
	t.extensionLookup = make(map[string]*ResolvedTool)
	t.resolved = make([]*ResolvedTool, 0)

	result := t.readToolOrder()
	t.initTools(result.Order, result.Found)
}

func (t *tools) initTools(tools []ToolKind, resultFoundInEnvVar bool) {
	t.extensionLookup = make(map[string]*ResolvedTool)
	t.pathLookup = make(map[string]*ResolvedTool, 0)
	t.resolved = make([]*ResolvedTool, 0)

	for _, tool := range t.sort(tools, resultFoundInEnvVar) {
		t.addToolWithSettings(string(tool.Kind), tool.Kind, tool.AutoRefresh, tool.IsMdi, tool.SupportsText, tool.RequiresTarget, tool.BinaryExtensions, &tool.Windows, &tool.Linux, &tool.Osx)
	}

	//add custom to the start
}

func (t *tools) sort(order []ToolKind, throwForNoTool bool) []*ToolDefinition {
	foundDefinitions := make([]*ToolDefinition, 0)
	allTools := make([]*ToolDefinition, len(AllDefinedTools))
	copy(allTools, AllDefinedTools)

	for _, k := range order {
		definition, found := t.findByKind(allTools, k)
		if !found {
			if !throwForNoTool {
				continue
			}
			panic(fmt.Sprintf("`DiffEngine_ToolOrder` is configured to use '%s' but it is not installed.", k))
		}

		foundDefinitions = append(foundDefinitions, definition)
		allTools = removeDefinition(allTools, definition)
	}

	for _, d := range allTools {
		foundDefinitions = append(foundDefinitions, d)
	}

	return foundDefinitions
}

func removeDefinition(slice []*ToolDefinition, e *ToolDefinition) []*ToolDefinition {
	index := getIndexOfDefinition(slice, e)
	if index != -1 {
		slice[index] = new(ToolDefinition)
		slice = append(slice[:index], slice[index+1:]...)
	}
	return slice
}

func getIndexOfDefinition(slice []*ToolDefinition, e *ToolDefinition) int {
	for k, v := range slice {
		if v == e {
			return k
		}
	}
	return -1 // not found.
}

func (t *tools) findByKind(definedTools []*ToolDefinition, kind ToolKind) (result *ToolDefinition, found bool) {
	for _, defined := range definedTools {
		if defined.Kind == kind {
			return defined, true
		}
	}
	return &ToolDefinition{}, false
}

func (t *tools) AddTool(name string, autoRefresh bool, isMdi bool, supportsText bool, requiresTarget bool,
	targetLeftArguments BuildArguments, targetRightArguments BuildArguments,
	exePath string, binaryExtensions []string) (*ResolvedTool, bool) {
	return t.addTool(name, None, autoRefresh, isMdi, supportsText, requiresTarget,
		binaryExtensions, exePath, targetLeftArguments, targetRightArguments)
}

func (t *tools) addToolWithSettings(name string, diffTool ToolKind, autoRefresh, isMdi, supportsText, requiresTarget bool,
	binaryExtensions []string, windows, linux, osx *OsSettings) (*ResolvedTool, bool) {

	if windows == nil && linux == nil && osx == nil {
		panic("must define settings for at least one OS.")
	}

	exe, left, right, found := resolveFromOsSettings(windows, linux, osx)
	if !found {
		return nil, false
	}

	return t.addTool(name, diffTool, autoRefresh, isMdi, supportsText,
		requiresTarget, binaryExtensions, exe, left, right)
}

func (t *tools) addTool(name string, diffTool ToolKind, autoRefresh bool, isMdi bool,
	supportsText bool, requiresTarget bool, binaries []string, exePath string,
	targetLeftArguments BuildArguments, targetRightArguments BuildArguments) (*ResolvedTool, bool) {

	guard.AgainstEmpty(name)
	if t.toolExists(name) {
		panic(fmt.Sprintf("Kind with Name already exists. Name: %s", name))
	}

	resolvedExePath, found := finder.TryFind(exePath)
	if !found {
		return nil, false
	}

	tool := &ResolvedTool{
		Name:             name,
		Kind:             diffTool,
		ExePath:          resolvedExePath,
		RightArguments:   targetRightArguments,
		LeftArguments:    targetLeftArguments,
		IsMdi:            isMdi,
		AutoRefresh:      autoRefresh,
		BinaryExtensions: binaries,
		RequiresTarget:   requiresTarget,
		SupportsText:     supportsText,
	}

	t.AddResolvedToolAtStart(tool)

	return tool, true
}

func (t *tools) toolExists(name string) bool {
	for _, resolve := range t.resolved {
		if resolve.Name == name {
			return true
		}
	}
	return false
}

type orderResult struct {
	Found bool
	Order []ToolKind
}

func (t *tools) readToolOrder() orderResult {
	diffOrder, found := os.LookupEnv("DiffEngine_ToolOrder")
	var order []ToolKind
	if found {
		order = t.parseEnvironment(diffOrder)
	} else {
		order = AllTools
	}

	return orderResult{
		Order: order,
		Found: found,
	}
}

func (t *tools) parseEnvironment(diffOrder string) []ToolKind {
	sep := func(r rune) bool {
		return r == ',' || r == '|' || r == ' '
	}

	tools := make([]ToolKind, 0)
	for _, toolString := range strings.FieldsFunc(diffOrder, sep) {
		tools = append(tools, ToolKind(toolString))
	}
	return tools
}

func (t *tools) AddResolvedToolAtStart(tool *ResolvedTool) {
	t.resolved = append([]*ResolvedTool{tool}, t.resolved...)
	for _, ext := range tool.BinaryExtensions {
		cleanedExtension := file.getFileExtension(ext)
		t.extensionLookup[cleanedExtension] = tool
	}
	t.pathLookup[tool.ExePath] = tool
}
