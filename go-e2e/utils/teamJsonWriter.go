package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const teamFileName = "team.json"

type TeamJson struct {
	RepositoriesPath *string   `json:"repositoriesPath,omitempty"`
	Repositories     *[]string `json:"repositories,omitempty"`
}

type TeamConfigOption func(*TeamJson)

func WriteTeamJsonTo(dir string, withOptions ...TeamConfigOption) {
	createDirError := os.MkdirAll(dir, 0755)
	if createDirError != nil {
		panic("failed to create directory " + dir + ": " + createDirError.Error())
	}

	config := makeTeamConfig(withOptions...)
	jsonBytes, marshallError := json.MarshalIndent(config, "", "  ")
	if marshallError != nil {
		panic("failed to marshal TeamJson to JSON: " + marshallError.Error())
	}

	filePath := filepath.Join(dir, teamFileName)

	writeError := os.WriteFile(filePath, jsonBytes, 0644)
	if writeError != nil {
		panic("failed to write TeamJson to file " + filePath + ": " + writeError.Error())
	}
}

func makeTeamConfig(options ...TeamConfigOption) *TeamJson {
	config := &TeamJson{}

	for _, option := range options {
		option(config)
	}

	return config
}

func WithRepositoriesPath(path string) TeamConfigOption {
	return func(c *TeamJson) {
		c.RepositoriesPath = &path
	}
}

func WithRepositories(repos []string) TeamConfigOption {
	return func(c *TeamJson) {
		c.Repositories = &repos
	}
}
