_config_file_name() {
	echo "config.json"
}

config_file_path() {
  local command_name="$1"

  bats_load_library 'commands/run/run_command_location.bash'
  echo "$(run_command_location)/$command_name/$(_config_file_name)"
}

_write_to_config_file() {
  local command_name="$1"
  local field_name="$2"
  local field_value="$3"

  write_json_field "$(config_file_path "$command_name")" "$field_name" "$field_value"
}

write_empty_json_object_as_config() {
  local command_name="$1"

  bats_load_library "json_writer.bash"
  write_empty_json_if_file_does_not_exist "$(config_file_path "$command_name")"
}

write_short_description() {
  local command_name="$1"
  local short_description="$2"

  bats_load_library "json_writer.bash"
  _write_to_config_file "$command_name" "shortDescription" "$(to_json_string "$short_description")"
}