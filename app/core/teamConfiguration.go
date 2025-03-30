package core

import (
	"errors"
	"github.com/spf13/viper"
)

const (
	teamFileName            = "team"
	teamFileExtension       = "json"
	defaultRepositoriesPath = "repositories"
	TeamFile                = teamFileName + "." + teamFileExtension
	RepositoriesPath        = "repositoriesPath"
)

var ErrCouldNotReadTeamFile = errors.New("could not read team file")

type TeamInfo struct {
	RepositoriesPath     string   `json:"repositoriesPath"`
	Repositories         []string `json:"repositories"`
	RepositoriesPrefixes []string `json:"repositoriesPrefixes"`
	CommitPrefixRegex    string   `json:"commitPrefixRegex"`
	BlockedBranches      []string `json:"blockedBranches"`
}

func LoadTeamConfiguration() (TeamInfo, error) {
	var teamInfo TeamInfo

	viper.AddConfigPath(GetExecutionPath())
	viper.SetConfigName(teamFileName)
	viper.SetConfigType(teamFileExtension)

	readErr := viper.ReadInConfig()
	if readErr != nil {
		return teamInfo, ErrCouldNotReadTeamFile
	}

	unmarshalErr := viper.Unmarshal(&teamInfo)
	if unmarshalErr != nil {
		return teamInfo, ErrCouldNotReadTeamFile
	}

	if teamInfo.RepositoriesPath == "" {
		teamInfo.RepositoriesPath = defaultRepositoriesPath
	}

	return teamInfo, nil
}
