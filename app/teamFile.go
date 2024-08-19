package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var teamFileName = "team.json"

type TeamInfo struct {
	RepositoriesPath string   `json:"repositoriesPath"`
	Repositories     []string `json:"repositories"`
}

func readTeamInfo() TeamInfo {
	fileContent, _ := os.Open(getExecutablePath() + "/" + teamFileName)

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var info TeamInfo
	jsonError := json.Unmarshal([]byte(byteResult), &info)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	return info
}
