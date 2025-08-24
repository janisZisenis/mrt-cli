_config_file_name() {
	echo "config.json"
}

config_file_path() {
  local commandName="$1"

  bats_load_library 'commands/run/run_command_location.bash'
  echo "$(run_command_location)/$commandName/$(_config_file_name)"
}

_write_to_config_file() {
  local commandName="$1"
  local fieldName="$2"
  local fieldValue="$3"

  write_json_field "$(config_file_path "$commandName")" "$fieldName" "$fieldValue"
}

write_empty_json_object_as_config() {
  local commandName="$1"

  bats_load_library "json_writer.bash"
  write_empty_json_if_file_does_not_exist "$(config_file_path "$commandName")"
}

write_short_description() {
  local commandName="$1"
  local shortDescription="$2"

  bats_load_library "json_writer.bash"
  _write_to_config_file "$commandName" "shortDescription" "$(to_json_string "$shortDescription")"
}