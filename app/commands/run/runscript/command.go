package runscript

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"mrt-cli/app/core"
	"mrt-cli/app/log"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // See GLOBAL_VIPER_STATE_FIX.md
var configMutex sync.Mutex

type CommandConfig struct {
	ShortDescription string `json:"shortDescription"`
}

func GetScriptsPath() string {
	return "/run/*/" + core.CommandFileName()
}

func LoadCommandConfig(commandPath string) CommandConfig {
	configMutex.Lock()
	defer configMutex.Unlock()

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
		//nolint:gocritic // See GLOBAL_VIPER_STATE_FIX.md
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
