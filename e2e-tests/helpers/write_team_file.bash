_teamFileName() {
	echo "team.json"
}

_teamFilePath() {
  bats_load_library 'fixtures/common_fixture.bash'
  echo "$(testEnvDir)/$(_teamFileName)"
}


_writeToTeamFile() {
  bats_load_library 'json_writer.bash'
  writeJsonField "$(_teamFilePath)" "$1" "$2"
}

writeBlockedBranches() {
  bats_load_library 'json_writer.bash'
	_writeToTeamFile "blockedBranches" "$(toJsonArray "$@")"
}

writeRepositoriesPrefixes() {
  bats_load_library 'json_writer.bash'
	_writeToTeamFile "repositoriesPrefixes" "$(toJsonArray "$@")"
}

writeRepositoriesUrls() {
  bats_load_library 'json_writer.bash'
	_writeToTeamFile "repositories" "$(toJsonArray "$@")"
}

writeRepositoriesPath() {
  bats_load_library 'json_writer.bash'
  _writeToTeamFile "repositoriesPath" "$(toJsonString "$1")"
}

writeCommitPrefixRegex() {
  bats_load_library 'json_writer.bash'
  _writeToTeamFile "commitPrefixRegex" "$(toJsonString "$1")"
}