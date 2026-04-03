package teamconfig

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

type Option func(*TeamJSON)

type Writer struct {
	dir string
}

func NewWriter(dir string) *Writer {
	return &Writer{dir: dir}
}

func (w *Writer) Write(withOptions ...Option) {
	createDirError := os.MkdirAll(w.dir, 0o750)
	if createDirError != nil {
		panic("failed to create directory " + w.dir + ": " + createDirError.Error())
	}

	config := makeConfig(withOptions...)
	jsonBytes, marshallError := json.MarshalIndent(config, "", "  ")
	if marshallError != nil {
		panic("failed to marshal TeamJSON to JSON: " + marshallError.Error())
	}

	filePath := filepath.Join(w.dir, teamFileName)

	writeError := os.WriteFile(filePath, jsonBytes, 0o600)
	if writeError != nil {
		panic("failed to write TeamJSON to file " + filePath + ": " + writeError.Error())
	}
}

func makeConfig(options ...Option) *TeamJSON {
	config := &TeamJSON{}

	for _, option := range options {
		option(config)
	}

	return config
}

func WithRepositoriesPath(path string) Option {
	return func(c *TeamJSON) {
		c.RepositoriesPath = &path
	}
}

func WithRepositories(repos []string) Option {
	return func(c *TeamJSON) {
		c.Repositories = &repos
	}
}

func WithRepositoriesPrefixes(prefixes []string) Option {
	return func(c *TeamJSON) {
		c.RepositoriesPrefixes = &prefixes
	}
}
