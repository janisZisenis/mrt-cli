package core

import (
	"errors"
	"github.com/spf13/viper"
)

const teamFileName = "team"
const teamFileExtension = "json"
const defaultRepositoriesPath = "repositories"
const TeamFile = teamFileName + "." + teamFileExtension
const RepositoriesPath = "repositoriesPath"

type TeamInfo struct {
	RepositoriesPath     string   `json:"repositoriesPath"`
	Repositories         []string `json:"repositories"`
	RepositoriesPrefixes []string `json:"repositoriesPrefixes"`
	CommitPrefixRegex    string   `json:"commitPrefixRegex"`
	BlockedBranches      []string `json:"blockedBranches"`
}

var CouldNotReadTeamFile = errors.New("could not read team file")

func LoadTeamConfiguration() (TeamInfo, error) {
	var teamInfo TeamInfo

	viper.AddConfigPath(GetExecutionPath())
	viper.SetConfigName(teamFileName)
	viper.SetConfigType(teamFileExtension)

	readErr := viper.ReadInConfig()
	unmarshalErr := viper.Unmarshal(&teamInfo)

	if teamInfo.RepositoriesPath == "" {
		teamInfo.RepositoriesPath = defaultRepositoriesPath
	}

	if readErr == nil && unmarshalErr == nil {
		return teamInfo, nil
	}

	return teamInfo, CouldNotReadTeamFile
}
