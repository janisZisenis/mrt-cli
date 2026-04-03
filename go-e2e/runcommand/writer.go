package runcommand

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	runDir         = "run"
	commandFile    = "command"
	configFile     = "config.json"
	dummyScript    = "#!/bin/bash\n"
	execPermission = 0o750
)

type Writer struct {
	teamDir string
}

func NewWriter(teamDir string) *Writer {
	return &Writer{teamDir: teamDir}
}

func (w *Writer) ConfigFilePath(commandName string) string {
	return filepath.Join(w.teamDir, runDir, commandName, configFile)
}

func (w *Writer) WriteDummyCommand(commandName string) {
	commandPath := filepath.Join(w.teamDir, runDir, commandName, commandFile)

	if err := os.MkdirAll(filepath.Dir(commandPath), 0o750); err != nil {
		panic("runcommand: failed to create command directory: " + err.Error())
	}

	if err := os.WriteFile(commandPath, []byte(dummyScript), execPermission); err != nil {
		panic("runcommand: failed to write command file: " + err.Error())
	}
}

func (w *Writer) WriteConfig(commandName string, options ...ConfigOption) {
	config := &commandConfig{}
	for _, opt := range options {
		opt(config)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic("runcommand: failed to marshal config: " + err.Error())
	}

	configPath := w.ConfigFilePath(commandName)

	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		panic("runcommand: failed to create config directory: " + err.Error())
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		panic("runcommand: failed to write config file: " + err.Error())
	}
}

func (w *Writer) WriteCorruptConfig(commandName string) {
	configPath := w.ConfigFilePath(commandName)

	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		panic("runcommand: failed to create config directory: " + err.Error())
	}

	if err := os.WriteFile(configPath, []byte{}, 0o600); err != nil {
		panic("runcommand: failed to write empty config file: " + err.Error())
	}
}

type commandConfig struct {
	ShortDescription *string `json:"shortDescription,omitempty"`
}

type ConfigOption func(*commandConfig)

func WithShortDescription(description string) ConfigOption {
	return func(c *commandConfig) {
		c.ShortDescription = &description
	}
}
