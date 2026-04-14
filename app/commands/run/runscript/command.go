package runscript

import (
	"errors"
	"mrt-cli/app/core"
	"mrt-cli/app/log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const (
	runScriptsDir          = "run"
	commandConfigFileName  = "config"
	commandConfigExtension = "json"
)

type CommandConfig struct {
	ShortDescription string `json:"shortDescription"`
}

func GetScriptsPath() string {
	return filepath.Join(runScriptsDir, "*", core.CommandFileName())
}

func LoadCommandConfig(commandPath string) CommandConfig {
	commandDir, _ := filepath.Abs(filepath.Dir(commandPath))
	setupDefaults(commandDir)

	var config CommandConfig

	viper.AddConfigPath(commandDir)
	viper.SetConfigName(commandConfigFileName)
	viper.SetConfigType(commandConfigExtension)

	readErr := viper.ReadInConfig()
	if readErr != nil {
		if errors.As(readErr, &viper.ConfigFileNotFoundError{}) {
			return defaultConfig(commandDir)
		}

		log.Errorf(
			"Error while reading %s/%s.%s",
			commandDir,
			commandConfigFileName,
			commandConfigExtension,
		)
		log.Errorf("%v", readErr)
		os.Exit(1)
	}

	_ = viper.Unmarshal(&config)

	return config
}

func setupDefaults(commandDir string) {
	viper.SetDefault("shortDescription", defaultConfig(commandDir).ShortDescription)
}

func defaultConfig(commandDir string) CommandConfig {
	description := "Executes run command " + filepath.Base(commandDir)
	return CommandConfig{
		ShortDescription: description,
	}
}

func MakeCommand(scriptName string, scriptPath string) *cobra.Command {
	config := LoadCommandConfig(scriptPath)

	command := &cobra.Command{
		Use:   scriptName,
		Short: config.ShortDescription,
		Run: func(_ *cobra.Command, args []string) {
			os.Exit(core.ExecuteScript(scriptPath, args))
		},
	}

	return command
}
