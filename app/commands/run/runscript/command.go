package runscript

import (
	"app/core"
	"app/log"
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type CommandConfig struct {
	ShortDescription string `json:"shortDescription"`
}

func GetScriptsPath() string {
	return "/run/*/" + core.CommandFileName()
}

func LoadCommandConfig(commandPath string) CommandConfig {
	commandDir := filepath.Dir(commandPath)
	setupDefaults(commandDir)

	var config CommandConfig

	viper.AddConfigPath(commandDir)
	configFileName := "config"
	configFileExtension := "json"
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)

	readErr := viper.ReadInConfig()
	if readErr != nil {
		if errors.As(readErr, &viper.ConfigFileNotFoundError{}) {
			return defaultConfig(commandDir)
		}

		log.Errorf("Error while reading %s/%s.%s", commandDir, configFileName, configFileExtension)
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

	var command = &cobra.Command{
		Use:   scriptName,
		Short: config.ShortDescription,
		Run: func(_ *cobra.Command, args []string) {
			command(scriptPath, args)
		},
	}

	return command
}

func command(scriptPath string, args []string) {
	scriptArgs := append([]string{core.GetAbsoluteExecutionPath()}, args...)
	exitCode := core.ExecuteScript(scriptPath, scriptArgs)
	os.Exit(exitCode)
}
