command_file_name() {
	echo "command"
}

write_dummy_command() {
  local commandLocation="$1"
  local commandName="$2"

  bats_load_library "scripts/script_writer.bash"
  write_dummy_script "$commandLocation/$commandName/$(command_file_name)"
}

write_stub_command() {
	local commandLocation="$1"
	local commandName="$2"
	local exitCode="$3"
	local output="$4"

  bats_load_library "scripts/script_writer.bash"
	write_stub_script "$commandLocation/$commandName/$(command_file_name)" "$exitCode" "$output"
}

write_spy_command() {
	local commandLocation="$1"
	local commandName="$2"

  bats_load_library "scripts/script_writer.bash"
	write_spy_script "$commandLocation/$commandName/$(command_file_name)"
}

assert_command_was_executed_with_parameters() {
	local commandLocation="$1"
	local commandName="$2"
	local expectedParameters="$3"

  bats_load_library "scripts/script_writer.bash"
	assert_script_was_executed_with_parameters "$commandLocation/$commandName/$(command_file_name)" "$expectedParameters"
}

assert_command_was_not_executed() {
	local commandLocation="$1"
	local commandName="$2"

  bats_load_library "scripts/script_writer.bash"
	assert_script_was_not_executed "$commandLocation/$commandName/$(command_file_name)"
}

write_std_err_command() {
	local commandLocation="$1"
	local commandName="$2"
	local error="$3"

  bats_load_library "scripts/script_writer.bash"
	write_std_err_script "$commandLocation/$commandName/$(command_file_name)" "$error"
}

write_command_requesting_input() {
	local commandLocation="$1"
	local commandName="$2"

  bats_load_library "scripts/script_writer.bash"
	write_script_requesting_input "$commandLocation/$commandName/$(command_file_name)"
}

assert_command_received_input() {
	local commandLocation="$1"
	local commandName="$2"
	local expectedInput="$3"

  bats_load_library "scripts/script_writer.bash"
  assert_script_received_input "$commandLocation/$commandName/$(command_file_name)" "$expectedInput"
}
