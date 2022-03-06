package diff

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"os"
	"strings"
)

type runner struct {
	Disabled   bool
	directory  string
	ciDetected CIDetected
	finder     *diffFinder
	tool       *tools
	counter    *instanceCounter
	proc       *processCleaner
	logger     Logger
}

const (
	envDiffEngineDisabled = "DiffEngine_Disabled"
)

type systemEnvReader struct{}

func (s *systemEnvReader) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Launch a new diff tool
func Launch(tempFile, targetFile string) LaunchResult {
	runner := newRunner(&systemEnvReader{})
	return runner.Launch(tempFile, targetFile)
}

// Kill the diff tool if it doesn't support MDI, is already running and has been
// opened to display a specific temp and target file.
func Kill(tempFile, targetFile string) {
	runner := newRunner(&systemEnvReader{})
	if runner.Disabled {
		return
	}

	extension := utils.File.GetFileExtension(tempFile)
	diffTool, found := runner.tool.tryFindForExtension(extension)
	if !found {
		runner.logger.Info("Extension not found. %s", extension)
		return
	}

	if diffTool.IsMdi {
		runner.logger.Info("DiffTool is Mdi so not killing. diffTool: %s", diffTool.ExePath)
		return
	}

	command := diffTool.buildCommand(tempFile, targetFile)
	runner.proc.Kill(command)
}

func newRunner(reader EnvReader) *runner {
	runner := &runner{
		ciDetected: checkCI(reader),
		finder:     newFinder(),
		tool:       newTools(),
		counter:    newInstanceCounter(reader),
		proc:       newProcessCleaner(),
		logger:     newLogger("runner"),
	}

	variable, found := reader.LookupEnv(envDiffEngineDisabled)
	if !found {
		variable = ""
	}

	if strings.ToLower(variable) == "true" || runner.ciDetected {
		runner.Disabled = true
	}

	return runner
}

// Launch runs a new diff tool that can handle the target file based on the file's extension.
func (r *runner) Launch(tempFile, targetFile string) LaunchResult {
	utils.Guard.GuardFiles(tempFile, targetFile)

	finder := func() (resolved *ResolvedTool, found bool) {
		extension := utils.File.GetFileExtension(tempFile)
		return r.tool.tryFindForExtension(extension)
	}

	return r.innerLaunch(finder, tempFile, targetFile)
}

// LaunchTool runs a specific diff tool
func (r *runner) LaunchTool(kind ToolKind, tempFile, targetFile string) LaunchResult {
	utils.Guard.GuardFiles(tempFile, targetFile)

	finder := func() (resolved *ResolvedTool, found bool) {
		return r.tool.tryFind(kind)
	}

	return r.innerLaunch(finder, tempFile, targetFile)
}

func (r *runner) innerLaunch(tryResolveTool TryResolveTool, tempFile, targetFile string) LaunchResult {
	tool, result, exit := r.ShouldExitLaunch(tryResolveTool, targetFile)
	if exit {
		//TODO: diff engine tray -> add move
		return result
	}

	args, cmd := tool.commandAndArguments(tempFile, targetFile)

	_, found := r.proc.GetProcessInfo(cmd)
	if found {
		if tool.AutoRefresh {
			//TODO: DiffEngineTray.AddMove
			return AlreadyRunningAndSupportsRefresh
		}

		r.KillIfMdi(tool, cmd)
	}

	if r.counter.ReachedMax() {
		//TODO: DiffEngineTray.AddMove
		return TooManyRunningDiffTools
	}

	_ = r.LaunchProcess(tool, args)
	//TODO: DiffEngineTray.AddMove

	return StartedNewInstance
}

// KillIfMdi kills the diff tool if it does not support MDI
func (r *runner) KillIfMdi(tool *ResolvedTool, command string) {
	if !tool.IsMdi {
		r.proc.Kill(command)
	}
}

// ShouldExitLaunch checks if the launched diff tool should be exitted
func (r *runner) ShouldExitLaunch(tryResolveTool TryResolveTool, targetFile string) (tool *ResolvedTool, result LaunchResult, exited bool) {
	if r.Disabled {
		return nil, Disabled, true
	}

	tool, found := tryResolveTool()
	if !found {
		return tool, NoDiffToolFound, true
	}

	if !r.tryCreate(tool, targetFile) {
		return tool, NoEmptyFileForExtension, true
	}

	return tool, NoLaunchResult, false
}

func (r *runner) tryCreate(tool *ResolvedTool, targetFile string) bool {
	targetExists := utils.File.Exists(targetFile)
	if tool.RequiresTarget && !targetExists {
		if !utils.File.TryCreateFile(targetFile, true) {
			return false
		}
	}
	return true
}

// LaunchProcess runs an external process with given arguments
func (r *runner) LaunchProcess(tool *ResolvedTool, arguments []string) int32 {

	out := make(chan runResult)
	go r.proc.RunCommand(out, tool.ExePath, arguments...)

	res := <-out
	if res.error != nil && res.pid == 0 {
		r.logger.Info("Failed to launch diff tool.\n" + tool.ExePath + " " + strings.Join(arguments, " "))
		return 0
	}

	return res.pid
}
