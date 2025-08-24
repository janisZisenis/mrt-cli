write_dummy_run_command() {
	local commandName="$1"

  bats_load_library 'commands/command_writer.bash'
  bats_load_library 'commands/run/run_command_location.bash'
  write_dummy_command "$(run_command_location)" "$commandName"
}

write_stub_run_command() {
	local commandName="$1"
	local exitCode="$2"
	local output="$3"

  bats_load_library 'commands/command_writer.bash'
  bats_load_library 'commands/run/run_command_location.bash'
  write_stub_command "$(run_command_location)" "$commandName" "$exitCode" "$output"
}

write_spy_run_command() {
	local commandName="$1"

  bats_load_library 'commands/command_writer.bash'
  bats_load_library 'commands/run/run_command_location.bash'
  write_spy_command "$(run_command_location)" "$commandName"
}

assert_run_command_was_executed_with_parameters() {
	local commandName="$1"
	local expectedParameters="$2"

  bats_load_library 'commands/command_writer.bash'
  bats_load_library 'commands/run/run_command_location.bash'
  assert_command_was_executed_with_parameters "$(run_command_location)" "$commandName" "$expectedParameters"
}

write_std_err_run_command() {
	local commandName="$1"
	local error="$2"

  bats_load_library 'commands/command_writer.bash'
	bats_load_library 'commands/run/run_command_location.bash'
	write_std_err_command "$(run_command_location)" "$commandName" "$error"
}

write_run_command_requesting_input() {
	local commandName="$1"

  bats_load_library 'commands/command_writer.bash'
	bats_load_library 'commands/run/run_command_location.bash'
	write_command_requesting_input "$(run_command_location)" "$commandName"
}

assert_run_command_received_input() {
	local commandName="$1"
	local expectedInput="$2"

  bats_load_library 'commands/command_writer.bash'
	bats_load_library 'commands/run/run_command_location.bash'
	assert_command_received_input "$(run_command_location)" "$commandName" "$expectedInput"
}