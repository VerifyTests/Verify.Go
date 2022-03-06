package diff

import (
	"bytes"
	"github.com/VerifyTests/Verify.Go/utils"
	"github.com/shirou/gopsutil/v3/process"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type processCleaner struct {
	commands []*processCommand
	logger   Logger
}

type processCommand struct {
	Command string
	Process int32
}

func newProcessCleaner() *processCleaner {
	cleaner := processCleaner{
		logger: newLogger("proc"),
	}
	cleaner.refresh()
	return &cleaner
}

func (p *processCleaner) refresh() {
	p.commands = p.findAllProcess()
}

func (p *processCleaner) tryTerminateProcess(pid int32) bool {
	proc, err := process.NewProcess(pid)
	if err == nil && proc != nil {
		running, err := proc.IsRunning()
		if running {
			err = proc.Terminate()
			return err == nil
		}
	}

	p.logger.Info("failed to kill process with id %d\n", pid)
	return false
}

func (p *processCleaner) findAllProcess() []*processCommand {
	processes := make([]*processCommand, 0)

	pids, _ := process.Pids()
	for _, pid := range pids {
		p, pe := process.NewProcess(pid)
		if pe != nil {
			continue
		}

		name, _ := p.Name()
		if len(name) == 0 {
			continue
		}

		cmdLine, _ := p.Cmdline()
		if len(cmdLine) == 0 {
			continue
		}

		proc := &processCommand{
			Process: p.Pid,
			Command: cmdLine,
		}
		processes = append(processes, proc)
	}

	return processes
}

func (p *processCleaner) GetProcessInfo(command string) (proc *processCommand, found bool) {
	utils.Guard.AgainstEmpty(command)
	if runtime.GOOS != "windows" {
		command = p.TrimCommand(command)
	}

	procCmd := p.findProcessByCommand(command)
	if procCmd != nil {
		p, err := process.NewProcess(procCmd.Process)
		if err == nil {
			running, err := p.IsRunning()
			if err == nil && running {
				return procCmd, true
			}
		}
	}
	return nil, false
}

func (p *processCleaner) TrimCommand(command string) string {
	return strings.ReplaceAll(command, "\"", "")
}

func (p *processCleaner) findProcessByCommand(command string) *processCommand {
	for _, p := range p.commands {
		if p.Command == command {
			return p
		}
	}
	return nil
}

type runResult struct {
	output []byte
	error  error
	pid    int32
}

func (p *processCleaner) RunCommand(ch chan<- runResult, name string, arg ...string) {
	cmd := exec.Command(name, arg...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	var err error
	var pid = int32(0)

	err = cmd.Start()
	if err == nil && cmd.Process != nil {
		pid = int32(cmd.Process.Pid)
	}

	err = cmd.Process.Release()
	if err != nil {
		return
	}

	time.Sleep(time.Second * 3)

	ch <- runResult{
		output: buf.Bytes(),
		error:  err,
		pid:    pid,
	}
}

func (p *processCleaner) Kill(command string) {
	utils.Guard.AgainstEmpty(command)

	if runtime.GOOS != "windows" {
		command = p.TrimCommand(command)
	}

	matchingCommands := p.findCommands(command)
	p.logger.Info("Kill: %s. Matching count: %d", command, len(matchingCommands))

	if len(matchingCommands) == 0 {
		p.logger.Info("No matching commands.")
	}

	for _, c := range matchingCommands {
		p.TerminateProcessIfExists(c.Process)
	}
}

func (p *processCleaner) getCommands() []string {
	matches := make([]string, 0)
	for _, c := range p.commands {
		matches = append(matches, c.Command)
	}
	return matches
}

func (p *processCleaner) findCommands(command string) []*processCommand {
	matches := make([]*processCommand, 0)
	for _, c := range p.commands {
		if strings.Compare(c.Command, command) == 0 {
			matches = append(matches, c)
		}
	}
	return matches
}

func (p *processCleaner) TerminateProcessIfExists(processId int32) {
	exited := p.tryTerminateProcess(processId)
	if exited {
		p.logger.Info("TerminateProcess. Id: %d.", processId)
	} else {
		p.logger.Info("Process not valid. Id: %d.", processId)
	}
}

func (p *processCleaner) IsRunning(pid int32) bool {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return false
	}

	running, err := proc.IsRunning()
	if err != nil {
		return false
	}

	return running
}

func (p *processCleaner) GetLineBreak() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
