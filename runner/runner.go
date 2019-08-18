package runner

import (
	"io"
	"os/exec"
)

func getRunCommand(debug bool) (string, []string) {
	if debug {
		return "dlv", []string{"exec",
			buildPath(),
			"--listen=:40000",
			"--headless",
			"--api-version=2",
			// "--log",
		}
	}

	return buildPath(), []string{}
}

func run(debug bool) bool {
	runnerLog("Running (isDebug: %v)...", debug)

	cname, args := getRunCommand(debug)
	cmd := exec.Command(cname, args...)

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
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
