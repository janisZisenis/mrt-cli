package fixtures

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	runCommandDir = "run"
	configFile    = "config.json"
)

type RunCommandFixture struct {
	*CommandFixture
}

func NewRunCommandFixture(repoDir string) *RunCommandFixture {
	return &RunCommandFixture{
		CommandFixture: NewCommandFixture(repoDir, runCommandDir),
	}
}

func (f *RunCommandFixture) ConfigFilePath(commandName string) string {
	return filepath.Join(filepath.Dir(f.CommandPath(commandName)), configFile)
}

func (f *RunCommandFixture) WriteConfig(commandName string, options ...runCommandConfigOption) {
	config := &runCommandConfig{}
	for _, opt := range options {
		opt(config)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic("runcommand: failed to marshal config: " + err.Error())
	}

	f.writeConfigFile(commandName, data)
}

func (f *RunCommandFixture) WriteCorruptConfig(commandName string) {
	f.writeConfigFile(commandName, []byte{})
}

func (f *RunCommandFixture) writeConfigFile(commandName string, data []byte) {
	configPath := f.ConfigFilePath(commandName)

	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		panic("runcommand: failed to create config directory: " + err.Error())
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		panic("runcommand: failed to write config file: " + err.Error())
	}
}

type runCommandConfig struct {
	ShortDescription *string `json:"shortDescription,omitempty"`
}

type runCommandConfigOption func(*runCommandConfig)

func WithShortDescription(description string) runCommandConfigOption {
	return func(c *runCommandConfig) {
		c.ShortDescription = &description
	}
}
