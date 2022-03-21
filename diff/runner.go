package diff

import (
	"github.com/VerifyTests/Verify.Go/tray"
	"github.com/VerifyTests/Verify.Go/utils"
	"os"
	"strings"
)

type runner struct {
	directory  string
	ciDetected CIDetected
	finder     *diffFinder
	tool       *Tools
	counter    *instanceCounter
	proc       *processCleaner
	tray       *tray.Client
	logger     Logger
	disabled   bool
}

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
	envReader := &systemEnvReader{}

	if checkDisabled(envReader) {
		return
	}

	runner := newRunner(envReader)

	extension := utils.File.GetFileExtension(tempFile)
	diffTool, found := runner.tool.TryFindForExtension(extension)
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
		disabled:   checkDisabled(reader),
		ciDetected: checkCI(reader),
		finder:     newDiffFinder(),
		tool:       NewTools(),
		counter:    newInstanceCounter(reader),
		proc:       newProcessCleaner(),
		tray:       tray.NewClient(),
		logger:     newLogger("runner"),
	}

	return runner
}

// Launch runs a new diff tool that can handle the target file based on the file's extension.
func (r *runner) Launch(tempFile, targetFile string) LaunchResult {
	utils.Guard.GuardFiles(tempFile, targetFile)

	finder := func() (resolved *ResolvedTool, found bool) {
		extension := utils.File.GetFileExtension(tempFile)
		return r.tool.TryFindForExtension(extension)
	}

	return r.innerLaunch(finder, tempFile, targetFile)
}

// LaunchTool runs a specific diff tool
func (r *runner) LaunchTool(kind ToolKind, tempFile, targetFile string) LaunchResult {
	utils.Guard.GuardFiles(tempFile, targetFile)

	finder := func() (resolved *ResolvedTool, found bool) {
		return r.tool.TryFind(kind)
	}

	return r.innerLaunch(finder, tempFile, targetFile)
}

func (r *runner) innerLaunch(tryResolveTool TryResolveTool, tempFile, targetFile string) LaunchResult {
	tool, result, exit := r.ShouldExitLaunch(tryResolveTool, targetFile)
	if exit {
		r.tray.AddMove(tempFile, targetFile, "", nil, false, -1)
		return result
	}

	args, cmd := tool.commandAndArguments(tempFile, targetFile)

	canKill := !tool.IsMdi
	processCommand, found := r.proc.GetProcessInfo(cmd)
	if found {
		if tool.AutoRefresh {
			r.tray.AddMove(tempFile, targetFile, tool.ExePath, args, canKill, processCommand.Process)
			return AlreadyRunningAndSupportsRefresh
		}

		r.KillIfMdi(tool, cmd)
	}

	if r.counter.ReachedMax() {
		r.tray.AddMove(tempFile, targetFile, tool.ExePath, args, canKill, 0)
		return TooManyRunningDiffTools
	}

	processId := r.LaunchProcess(tool, args)
	r.tray.AddMove(tempFile, targetFile, tool.ExePath, args, canKill, processId)

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
	if r.disabled {
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
