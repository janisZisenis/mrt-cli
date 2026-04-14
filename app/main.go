package main

import (
	"mrt-cli/app/commands/githook"
	"mrt-cli/app/commands/run"
	"mrt-cli/app/commands/setup"
	"mrt-cli/app/commands/version"
	"mrt-cli/app/core"
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
	core.SetTeamDirectory(readTeamDir())

	rootCmd := &cobra.Command{Use: filepath.Base(os.Args[0])}

	rootCmd.AddCommand(setup.MakeCommand())
	rootCmd.AddCommand(githook.MakeCommand())
	rootCmd.AddCommand(run.MakeCommand())
	rootCmd.AddCommand(version.MakeCommand(semver, commit, date))

	rootCmd.PersistentFlags().String("team-dir", "", "Specifies the path to the team directory.")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
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
