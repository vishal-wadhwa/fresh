package runner

import (
	"context"
	"io"
	"os/exec"
	"strconv"
)

func getDebugCommand(pid int) (string, []string) {
	return "dlv", []string{"attach",
		strconv.Itoa(pid),
		"--listen=:40000",
		"--headless",
		"--api-version=2",
		"--accept-multiclient"}
	// "--log",
}

// KillFn specifies the type for function which kills processes spawned by fresh-dlv
type KillFn func()

func run(debug bool) KillFn {
	runnerLog("Running (isDebug: %v)...", debug)
	ctx, cancel := context.WithCancel(context.Background())

	var debugCmd *exec.Cmd
	var cancelDebugger context.CancelFunc
	cmd := exec.CommandContext(ctx, buildPath())

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

	writer := appLogWriter{}
	if debug {
		writer := debuggerLogWriter{}
		cname, args := getDebugCommand(cmd.Process.Pid)
		ctx, cancelDebugger = context.WithCancel(context.Background())
		debugCmd = exec.CommandContext(ctx, cname, args...)

		stdout, err := debugCmd.StdoutPipe()
		if err != nil {
			fatal(err)
		}

		stderr, err := debugCmd.StderrPipe()
		if err != nil {
			fatal(err)
		}

		if err := debugCmd.Start(); err != nil {
			fatal(err)
		}

		go io.Copy(writer, stderr)
		go io.Copy(writer, stdout)
	}
	go io.Copy(writer, stderr)
	go io.Copy(writer, stdout)

	return func() {
		cancel()
		if cancelDebugger != nil {
			cancelDebugger()
		}
	}
}
