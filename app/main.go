package main

import (
	"app/commands/githook"
	"app/commands/setup"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func main() {
	var rootCmd = &cobra.Command{Use: filepath.Base(os.Args[0])}
	rootCmd.AddCommand(setup.MakeCommand())
	rootCmd.AddCommand(githook.MakeCommand())
	rootCmd.Execute()
}
