package runner

import (
	"io"
	"os/exec"
	"syscall"
)

func getRunCommand(debug bool) (string, []string) {
	if debug {
		return "dlv", []string{"exec",
			buildPath(),
			"--continue",
			"--listen=:40000",
			"--headless",
			"--api-version=2",
			"--accept-multiclient",
			// "--log",
		}
	}

	return buildPath(), []string{}
}

func run(debug bool) bool {
	runnerLog("Running (isDebug: %v)...", debug)

	cname, args := getRunCommand(debug)
	cmd := exec.Command(cname, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	var writer io.Writer = appLogWriter{}
	if debug {
		writer = debuggerLogWriter{}
	}
	go io.Copy(writer, stderr)
	go io.Copy(writer, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		pgid, err := syscall.Getpgid(pid)
		if err != nil {
			runnerLog("error getting pgid for process %d - %s", pid, err.Error())
		}
		runnerLog("Killing pgid(pid) - %d(%d)", pgid, pid)
		err = syscall.Kill(-pid, syscall.SIGINT)
		if err != nil {
			runnerLog("could not kill pid - %d", pid)
		}
	}()

	return true
}
