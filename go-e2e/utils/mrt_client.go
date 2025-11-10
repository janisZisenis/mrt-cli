package utils

import (
	"os"
	"os/exec"
	"strings"
)

type Mrt struct {
	binaryName string
	command    *exec.Cmd
}

func MakeMrtCommand(binaryPath string, env []string) *Mrt {
	command := exec.Command(binaryPath)
	command.Env = append(os.Environ(), env...)

	return &Mrt{
		binaryName: binaryPath,
		command:    command,
	}
}

func (mrt *Mrt) RunInDirectory(directory string) *Mrt {
	mrt.command.Args = append(mrt.command.Args, "--team-dir", directory)

	return mrt
}

func (mrt *Mrt) Setup() *Mrt {
	mrt.command.Args = append(mrt.command.Args, "setup")

	return mrt
}

func (mrt *Mrt) Clone() *Mrt {
	mrt.command.Args = append(mrt.command.Args, "clone-repositories")

	return mrt
}

func (mrt *Mrt) Execute() *Output {
	byteOutput, err := mrt.command.CombinedOutput()
	output := string(byteOutput)

	if err != nil {
		panic("executing mrt command failed: " + output)
	}

	return MakeOutput(SplitLines(output))
}

func SplitLines(output string) []string {
	if output == "" {
		return []string{}
	}

	return strings.Split(strings.TrimSpace(output), "\n")
}
