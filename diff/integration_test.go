//go:build integration
// +build integration

package diff

import (
	"path/filepath"
	"testing"
	"time"
)

func TestFindingProcessByName_Integration(t *testing.T) {
	env := newTestReader()
	env.lookup[envDiffEngineDisabled] = "false"

	temp := filepath.Join("../_testdata", "temp.txt")
	tempPath, _ := filepath.Abs(temp)
	target := filepath.Join("../_testdata", "target.txt")
	targetPath, _ := filepath.Abs(target)

	file.writeText(temp, "temp file")
	file.writeText(target, "target file")

	r := newRunner(env)
	vs, _ := r.tool.tryFind(VisualStudioCode)
	r.proc.Kill(vs.ExePath)

	time.Sleep(time.Second * 3)

	result := r.LaunchTool(VisualStudioCode, tempPath, targetPath)

	if result != StartedNewInstance {
		t.Fatalf("should start a new instance")
	}
}

func TestToolsInitialization_Integration(t *testing.T) {
	tool := tools{}
	tool.initTools([]ToolKind{VisualStudioCode}, false)

	if len(tool.resolved) != 2 {
		t.Fatalf("should find two diff tools")
	}
}
