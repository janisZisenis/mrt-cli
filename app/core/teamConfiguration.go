package core

import (
	"app/log"
	"github.com/spf13/viper"
	"os"
)

var teamFileName = "team"
var teamFileExtension = "json"
var TeamFile = teamFileName + "." + teamFileExtension
var RepositoriesPath = "repositoriesPath"

type TeamInfo struct {
	RepositoriesPath     string   `json:"repositoriesPath"`
	Repositories         []string `json:"repositories"`
	RepositoriesPrefixes []string `json:"repositoriesPrefixes"`
	CommitPrefixRegex    string   `json:"commitPrefixRegex"`
	BlockedBranches      []string `json:"blockedBranches"`
}

func LoadTeamConfiguration() TeamInfo {
	var teamInfo TeamInfo

	viper.AddConfigPath(GetExecutablePath())
	viper.SetConfigName(teamFileName)
	viper.SetConfigType(teamFileExtension)

	readErr := viper.ReadInConfig()
	unmarshalErr := viper.Unmarshal(&teamInfo)

	if readErr != nil || unmarshalErr != nil {
		log.Error("Could not read team file. Please make sure a \"" + TeamFile + "\" file exists next " +
			"to the executable and that it follows proper JSON syntax")
		os.Exit(1)
	}

	if teamInfo.RepositoriesPath == "" {
		teamInfo.RepositoriesPath = "repositories"
	}

	return teamInfo
}
