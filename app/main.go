package main

import (
	"app/commands/githook"
	"app/commands/run"
	"app/commands/setup"
	"app/commands/version"
	"app/core"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	semver = "v0.0.0"
	commit = "000000"
	date   = "1979-01-01"
)

func main() {
	args := os.Args[1:]
	for i, arg := range args {
		core.TeamDirectory = core.GetExecutionPath()
		if strings.HasPrefix(arg, "--team-dir=") {
			core.TeamDirectory = strings.SplitN(arg, "=", 2)[1]
		} else if arg == "--team-dir" && i+1 < len(args) {
			core.TeamDirectory = args[i+1]
		}
	}

	var rootCmd = &cobra.Command{Use: filepath.Base(os.Args[0])}

	rootCmd.AddCommand(setup.MakeCommand(core.TeamDirectory))
	rootCmd.AddCommand(githook.MakeCommand())
	rootCmd.AddCommand(run.MakeCommand(core.TeamDirectory))
	rootCmd.AddCommand(version.MakeCommand(semver, commit, date))

	rootCmd.PersistentFlags().StringVar(&core.TeamDirectory, "team-dir", "", "Specifies the path to the team directory.")
	_ = rootCmd.Execute()
}
