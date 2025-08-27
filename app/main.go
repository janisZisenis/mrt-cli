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

//nolint:gochecknoglobals // These variables are injected at build time via -ldflags, it's not really global variables.
var (
	semver = "v0.0.0"
	commit = "000000"
	date   = "1979-01-01"
)

func main() {
	teamDirFromFlag := readTeamDir()
	core.SetTeamDirectory(teamDirFromFlag)

	executionPath := core.GetExecutionPath()

	var rootCmd = &cobra.Command{Use: filepath.Base(os.Args[0])}

	rootCmd.AddCommand(setup.MakeCommand(executionPath))
	rootCmd.AddCommand(githook.MakeCommand())
	rootCmd.AddCommand(run.MakeCommand(executionPath))
	rootCmd.AddCommand(version.MakeCommand(semver, commit, date))

	rootCmd.PersistentFlags().StringVar(&executionPath, "team-dir", "", "Specifies the path to the team directory.")
	_ = rootCmd.Execute()
}

func readTeamDir() *string {
	args := os.Args[1:]
	for i, arg := range args {
		_, after, found := strings.Cut(arg, "--team-dir=")

		if found {
			return &after
		} else if arg == "--team-dir" && i+1 < len(args) {
			return &args[i+1]
		}
	}

	return nil
}
