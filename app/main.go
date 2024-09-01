package main

import (
	"app/commands/setup"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func main() {
	var rootCmd = &cobra.Command{Use: filepath.Base(os.Args[0])}
	rootCmd.AddCommand(setup.Command)
	rootCmd.Execute()
}
