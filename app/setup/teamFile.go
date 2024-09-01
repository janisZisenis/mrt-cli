package setup

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var TeamFileName = "team.json"

type TeamInfo struct {
	RepositoriesPath     string   `json:"repositoriesPath"`
	Repositories         []string `json:"repositories"`
	RepositoriesPrefixes []string `json:"repositoriesPrefixes"`
}

func ReadTeamInfo() TeamInfo {
	fileContent, _ := os.Open(GetExecutablePath() + "/" + TeamFileName)

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var info TeamInfo
	jsonError := json.Unmarshal([]byte(byteResult), &info)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	return info
}
