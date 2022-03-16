package diff

import (
	"fmt"
	"github.com/VerifyTests/Verify.Go/utils"
	"os"
	"strings"
)

type Tools struct {
	resolved        []*ResolvedTool
	pathLookup      map[string]*ResolvedTool
	extensionLookup map[string]*ResolvedTool
}

// NewTools creates a new diff tool accessor
func NewTools() *Tools {
	t := &Tools{}
	t.reset()
	return t
}

// TryFind finds a tool by the provided kind
func (t *Tools) TryFind(kind ToolKind) (tool *ResolvedTool, found bool) {
	for _, rt := range t.resolved {
		if rt.Kind == kind {
			return rt, true
		}
	}
	return nil, false
}

// TryFindByPath finds a tool by the provided path
func (t *Tools) TryFindByPath(path string) (tool *ResolvedTool, found bool) {
	if tool, ok := t.pathLookup[path]; ok {
		return tool, true
	}
	return nil, false
}

//TryFindForExtension finds a tool based on the provided extension
func (t *Tools) TryFindForExtension(extension string) (tool *ResolvedTool, found bool) {
	extension = utils.File.GetFileExtension(extension)
	if utils.File.IsText(extension) {
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

func (t *Tools) reset() {
	t.extensionLookup = make(map[string]*ResolvedTool)
	t.resolved = make([]*ResolvedTool, 0)

	result := t.readToolOrder()
	t.initTools(result.Order, result.Found)
}

func (t *Tools) initTools(tools []ToolKind, resultFoundInEnvVar bool) {
	t.extensionLookup = make(map[string]*ResolvedTool)
	t.pathLookup = make(map[string]*ResolvedTool, 0)
	t.resolved = make([]*ResolvedTool, 0)

	for _, tool := range t.sort(tools, resultFoundInEnvVar) {
		t.addToolWithSettings(string(tool.Kind), tool.Kind, tool.AutoRefresh, tool.IsMdi, tool.SupportsText, tool.RequiresTarget, tool.BinaryExtensions, &tool.Windows, &tool.Linux, &tool.Osx)
	}

	//add custom to the start
}

func (t *Tools) sort(order []ToolKind, throwForNoTool bool) []*ToolDefinition {
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

func (t *Tools) findByKind(definedTools []*ToolDefinition, kind ToolKind) (result *ToolDefinition, found bool) {
	for _, defined := range definedTools {
		if defined.Kind == kind {
			return defined, true
		}
	}
	return &ToolDefinition{}, false
}

func (t *Tools) AddTool(name string, autoRefresh bool, isMdi bool, supportsText bool, requiresTarget bool,
	targetLeftArguments BuildArguments, targetRightArguments BuildArguments,
	exePath string, binaryExtensions []string) (*ResolvedTool, bool) {
	return t.addTool(name, None, autoRefresh, isMdi, supportsText, requiresTarget,
		binaryExtensions, exePath, targetLeftArguments, targetRightArguments)
}

func (t *Tools) addToolWithSettings(name string, diffTool ToolKind, autoRefresh, isMdi, supportsText, requiresTarget bool,
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

func (t *Tools) addTool(name string, diffTool ToolKind, autoRefresh bool, isMdi bool,
	supportsText bool, requiresTarget bool, binaries []string, exePath string,
	targetLeftArguments BuildArguments, targetRightArguments BuildArguments) (*ResolvedTool, bool) {

	utils.Guard.AgainstEmpty(name)
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

func (t *Tools) toolExists(name string) bool {
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

func (t *Tools) readToolOrder() orderResult {
	diffOrder, found := os.LookupEnv("DiffEngine_ToolOrder")
	var order []ToolKind
	if found {
		order = t.parseEnvironment(diffOrder)
	} else {
		order = allTools
	}

	return orderResult{
		Order: order,
		Found: found,
	}
}

func (t *Tools) parseEnvironment(diffOrder string) []ToolKind {
	sep := func(r rune) bool {
		return r == ',' || r == '|' || r == ' '
	}

	tools := make([]ToolKind, 0)
	for _, toolString := range strings.FieldsFunc(diffOrder, sep) {
		tools = append(tools, ToolKind(toolString))
	}
	return tools
}

func (t *Tools) AddResolvedToolAtStart(tool *ResolvedTool) {
	t.resolved = append([]*ResolvedTool{tool}, t.resolved...)
	for _, ext := range tool.BinaryExtensions {
		cleanedExtension := utils.File.GetFileExtension(ext)
		t.extensionLookup[cleanedExtension] = tool
	}
	t.pathLookup[tool.ExePath] = tool
}
