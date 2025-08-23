_configFileName() {
	echo "config.json"
}

configFilePath() {
  local commandName="$1"

  bats_load_library 'commands/run/runCommandLocation.bash'
  echo "$(runCommandLocation)/$commandName/$(_configFileName)"
}

_writeToConfigFile() {
  local commandName="$1"
  local fieldName="$2"
  local fieldValue="$3"

  writeJsonField "$(configFilePath "$commandName")" "$fieldName" "$fieldValue"
}

writeEmptyJsonObjectAsConfig() {
  local commandName="$1"

  bats_load_library "jsonWriter.bash"
  writeEmptyJsonIfFileDoesNotExist "$(configFilePath "$commandName")"
}

writeShortDescription() {
  local commandName="$1"
  local shortDescription="$2"

  bats_load_library "jsonWriter.bash"
  _writeToConfigFile "$commandName" "shortDescription" "$(toJsonString "$shortDescription")"
}