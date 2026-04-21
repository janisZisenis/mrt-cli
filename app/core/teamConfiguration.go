package core

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	teamFileName            = "team"
	teamFileExtension       = "json"
	defaultRepositoriesPath = "repositories"
	TeamFile                = teamFileName + "." + teamFileExtension
	RepositoriesPath        = "repositoriesPath"
)

type TeamInfo struct {
	RepositoriesPath     string   `json:"repositoriesPath"`
	Repositories         []string `json:"repositories"`
	RepositoriesPrefixes []string `json:"repositoriesPrefixes"`
	CommitPrefixRegex    string   `json:"commitPrefixRegex"`
	BlockedBranches      []string `json:"blockedBranches"`
}

var (
	ErrCouldNotReadTeamFile    = errors.New("could not read team file")
	ErrInvalidRepositoriesPath = errors.New("repositoriesPath must be a relative path within the team repository")
)

func LoadTeamConfiguration(teamDir string) (TeamInfo, error) {
	var teamInfo TeamInfo

	viper.AddConfigPath(teamDir)
	viper.SetConfigName(teamFileName)
	viper.SetConfigType(teamFileExtension)

	readErr := viper.ReadInConfig()
	unmarshalErr := viper.Unmarshal(&teamInfo)

	if teamInfo.RepositoriesPath == "" {
		teamInfo.RepositoriesPath = defaultRepositoriesPath
	}

	if readErr == nil && unmarshalErr == nil {
		absTeamDir, err := filepath.Abs(teamDir)
		if err != nil {
			return teamInfo, ErrCouldNotReadTeamFile
		}

		if filepath.IsAbs(teamInfo.RepositoriesPath) {
			return teamInfo, ErrInvalidRepositoriesPath
		}

		resolved := filepath.Clean(filepath.Join(absTeamDir, teamInfo.RepositoriesPath))
		if !strings.HasPrefix(resolved+string(filepath.Separator), absTeamDir+string(filepath.Separator)) {
			return teamInfo, ErrInvalidRepositoriesPath
		}

		return teamInfo, nil
	}

	return teamInfo, ErrCouldNotReadTeamFile
}
