package core

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

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
	viper.SetConfigName("team")
	viper.SetConfigType("json")

	readErr := viper.ReadInConfig()
	unmarshalErr := viper.Unmarshal(&teamInfo)

	if readErr != nil || unmarshalErr != nil {
		fmt.Println("Could not read team file. Please make sure a \"team.json\" file exists next " +
			"to the executable and that it follows proper JSON syntax")
		os.Exit(1)
	}

	if teamInfo.RepositoriesPath == "" {
		teamInfo.RepositoriesPath = "repositories"
	}

	return teamInfo
}
