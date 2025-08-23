_teamFileName() {
	echo "team.json"
}

_teamFilePath() {
  bats_load_library 'fixtures/common_fixture.bash'
  echo "$(testEnvDir)/$(_teamFileName)"
}


_writeToTeamFile() {
  bats_load_library 'jsonWriter.bash'
  writeJsonField "$(_teamFilePath)" "$1" "$2"
}

writeBlockedBranches() {
  bats_load_library 'jsonWriter.bash'
	_writeToTeamFile "blockedBranches" "$(toJsonArray "$@")"
}

writeRepositoriesPrefixes() {
  bats_load_library 'jsonWriter.bash'
	_writeToTeamFile "repositoriesPrefixes" "$(toJsonArray "$@")"
}

writeRepositoriesUrls() {
  bats_load_library 'jsonWriter.bash'
	_writeToTeamFile "repositories" "$(toJsonArray "$@")"
}

writeRepositoriesPath() {
  bats_load_library 'jsonWriter.bash'
  _writeToTeamFile "repositoriesPath" "$(toJsonString "$1")"
}

writeCommitPrefixRegex() {
  bats_load_library 'jsonWriter.bash'
  _writeToTeamFile "commitPrefixRegex" "$(toJsonString "$1")"
}