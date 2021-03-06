package diff

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestFindingProcessByName_Integration(t *testing.T) {
	runTests, _ := strconv.ParseBool(os.Getenv("RUN_INTEGRATION_TESTS"))
	if !runTests {
		t.Skip("Skipping integration tests")
	}

	env := newTestReader()
	env.lookup[envDiffEngineDisabled] = "false"

	temp := filepath.Join("../_testdata", "temp.txt")
	tempPath, _ := filepath.Abs(temp)
	target := filepath.Join("../_testdata", "target.txt")
	targetPath, _ := filepath.Abs(target)

	utils.File.WriteText(temp, "temp file")
	utils.File.WriteText(target, "target file")

	r := newRunner(env)
	vs, _ := r.tool.TryFind(VisualStudioCode)
	r.proc.Kill(vs.ExePath)

	time.Sleep(time.Second * 3)

	result := r.LaunchTool(VisualStudioCode, tempPath, targetPath)

	if result != StartedNewInstance {
		t.Fatalf("should start a new instance")
	}
}

func TestToolsInitialization_Integration(t *testing.T) {
	runTests, _ := strconv.ParseBool(os.Getenv("RUN_INTEGRATION_TESTS"))
	if !runTests {
		t.Skip("Skipping integration tests")
	}

	tool := Tools{}
	tool.initTools([]ToolKind{VisualStudioCode}, false)

	if len(tool.resolved) != 2 {
		t.Fatalf("should find two diff tools")
	}
}
