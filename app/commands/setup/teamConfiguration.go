package setup

import (
	"fmt"
	"github.com/spf13/viper"
)

type TeamInfo struct {
	RepositoriesPath     string   `json:"repositoriesPath"`
	Repositories         []string `json:"repositories"`
	RepositoriesPrefixes []string `json:"repositoriesPrefixes"`
	BlockedBranches      []string `json:"blockedBranches"`
}

func LoadTeamConfiguration() (*TeamInfo, error) {
	var teamInfo *TeamInfo

	viper.AddConfigPath(GetExecutablePath())
	viper.SetConfigName("team")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not read team file: %s \n", err)
	}

	err = viper.Unmarshal(&teamInfo)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshall the team file: %s \n", err)
	}

	if teamInfo.RepositoriesPath == "" {
		teamInfo.RepositoriesPath = "repositories"
	}

	return teamInfo, nil
}
