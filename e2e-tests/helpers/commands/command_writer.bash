command_file_name() {
	echo "command"
}

write_dummy_command() {
  local command_location="$1"
  local command_name="$2"

  bats_load_library "scripts/script_writer.bash"
  write_dummy_script "$command_location/$command_name/$(command_file_name)"
}

write_stub_command() {
	local command_location="$1"
	local command_name="$2"
	local exit_code="$3"
	local output="$4"

  bats_load_library "scripts/script_writer.bash"
	write_stub_script "$command_location/$command_name/$(command_file_name)" "$exit_code" "$output"
}

write_spy_command() {
	local command_location="$1"
	local command_name="$2"

  bats_load_library "scripts/script_writer.bash"
	write_spy_script "$command_location/$command_name/$(command_file_name)"
}

assert_command_was_executed_with_parameters() {
	local command_location="$1"
	local command_name="$2"
	local expected_parameters="$3"

  bats_load_library "scripts/script_writer.bash"
	assert_script_was_executed_with_parameters "$command_location/$command_name/$(command_file_name)" "$expected_parameters"
}

assert_command_was_not_executed() {
	local command_location="$1"
	local command_name="$2"

  bats_load_library "scripts/script_writer.bash"
	assert_script_was_not_executed "$command_location/$command_name/$(command_file_name)"
}

write_std_err_command() {
	local command_location="$1"
	local command_name="$2"
	local error="$3"

  bats_load_library "scripts/script_writer.bash"
	write_std_err_script "$command_location/$command_name/$(command_file_name)" "$error"
}

write_command_requesting_input() {
	local command_location="$1"
	local command_name="$2"

  bats_load_library "scripts/script_writer.bash"
	write_script_requesting_input "$command_location/$command_name/$(command_file_name)"
}

assert_command_received_input() {
	local command_location="$1"
	local command_name="$2"
	local expected_input="$3"

  bats_load_library "scripts/script_writer.bash"
  assert_script_received_input "$command_location/$command_name/$(command_file_name)" "$expected_input"
}
