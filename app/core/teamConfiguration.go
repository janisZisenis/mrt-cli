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
	RepositoriesPath     string
	Repositories         []string
	RepositoriesPrefixes []string
	CommitPrefixRegex    string
	BlockedBranches      []string
}

type TeamJson struct {
	RepositoriesPath     string   `json:"repositoriesPath"`
	Repositories         []string `json:"repositories"`
	RepositoriesPrefixes []string `json:"repositoriesPrefixes"`
	CommitPrefixRegex    string   `json:"commitPrefixRegex"`
	BlockedBranches      []string `json:"blockedBranches"`
}

var CouldNotReadTeamFile = errors.New("could not read team file")

func LoadTeamConfiguration() (TeamInfo, error) {
	var teamJson TeamJson

	viper.AddConfigPath(GetExecutablePath())
	viper.SetConfigName(teamFileName)
	viper.SetConfigType(teamFileExtension)

	readErr := viper.ReadInConfig()
	unmarshalErr := viper.Unmarshal(&teamJson)

	if readErr != nil || unmarshalErr != nil {
		return makeEmptyTeamInfo(), CouldNotReadTeamFile
	}

	var teamInfo = makeTeamInfo(teamJson)
	if teamInfo.RepositoriesPath == "" {
		teamInfo.RepositoriesPath = defaultRepositoriesPath
	}

	return teamInfo, nil
}

func makeEmptyTeamInfo() TeamInfo {
	return TeamInfo{}
}

func makeTeamInfo(json TeamJson) TeamInfo {
	return TeamInfo{
		RepositoriesPath:     json.RepositoriesPath,
		Repositories:         json.Repositories,
		RepositoriesPrefixes: json.RepositoriesPrefixes,
		CommitPrefixRegex:    json.CommitPrefixRegex,
		BlockedBranches:      json.BlockedBranches,
	}
}
