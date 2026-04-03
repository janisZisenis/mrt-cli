package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const teamFileName = "team.json"

type TeamJSON struct {
	RepositoriesPath     *string   `json:"repositoriesPath,omitempty"`
	Repositories         *[]string `json:"repositories,omitempty"`
	RepositoriesPrefixes *[]string `json:"repositoriesPrefixes,omitempty"`
}

type TeamConfigOption func(*TeamJSON)

func WriteTeamJSONTo(dir string, withOptions ...TeamConfigOption) {
	createDirError := os.MkdirAll(dir, 0o750)
	if createDirError != nil {
		panic("failed to create directory " + dir + ": " + createDirError.Error())
	}

	config := makeTeamConfig(withOptions...)
	jsonBytes, marshallError := json.MarshalIndent(config, "", "  ")
	if marshallError != nil {
		panic("failed to marshal TeamJSON to JSON: " + marshallError.Error())
	}

	filePath := filepath.Join(dir, teamFileName)

	writeError := os.WriteFile(filePath, jsonBytes, 0o600)
	if writeError != nil {
		panic("failed to write TeamJSON to file " + filePath + ": " + writeError.Error())
	}
}

func makeTeamConfig(options ...TeamConfigOption) *TeamJSON {
	config := &TeamJSON{}

	for _, option := range options {
		option(config)
	}

	return config
}

func WithRepositoriesPath(path string) TeamConfigOption {
	return func(c *TeamJSON) {
		c.RepositoriesPath = &path
	}
}

func WithRepositories(repos []string) TeamConfigOption {
	return func(c *TeamJSON) {
		c.Repositories = &repos
	}
}

func WithRepositoriesPrefixes(prefixes []string) TeamConfigOption {
	return func(c *TeamJSON) {
		c.RepositoriesPrefixes = &prefixes
	}
}
