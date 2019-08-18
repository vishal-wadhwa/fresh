package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func getBuildCommand(debug bool) (string, []string) {
	args := []string{"build"}
	if debug {
		args = append(args, "-gcflags=\"all=-N -l\"")
	}
	buildDest := []string{"-o", buildPath(), root()}
	args = append(args, buildDest...)

	return "go", args
}

func build(debug bool) (string, bool) {
	buildLog("Building (isDebug: %v)...", debug)

	gocmd, args := getBuildCommand(debug)
	cmd := exec.Command(gocmd, args...)

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

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	return "", true
}
