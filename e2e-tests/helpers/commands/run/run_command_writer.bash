write_dummy_run_command() {
	local command_name="$1"

  bats_load_library 'commands/command_writer.bash'
  bats_load_library 'commands/run/run_command_location.bash'
  write_dummy_command "$(run_command_location)" "$command_name"
}

write_stub_run_command() {
	local command_name="$1"
	local exit_code="$2"
	local output="$3"

  bats_load_library 'commands/command_writer.bash'
  bats_load_library 'commands/run/run_command_location.bash'
  write_stub_command "$(run_command_location)" "$command_name" "$exit_code" "$output"
}

write_spy_run_command() {
	local command_name="$1"

  bats_load_library 'commands/command_writer.bash'
  bats_load_library 'commands/run/run_command_location.bash'
  write_spy_command "$(run_command_location)" "$command_name"
}

assert_run_command_was_executed_with_parameters() {
	local command_name="$1"
	local expected_parameters="$2"

  bats_load_library 'commands/command_writer.bash'
  bats_load_library 'commands/run/run_command_location.bash'
  assert_command_was_executed_with_parameters "$(run_command_location)" "$command_name" "$expected_parameters"
}

write_std_err_run_command() {
	local command_name="$1"
	local error="$2"

  bats_load_library 'commands/command_writer.bash'
	bats_load_library 'commands/run/run_command_location.bash'
	write_std_err_command "$(run_command_location)" "$command_name" "$error"
}

write_run_command_requesting_input() {
	local command_name="$1"

  bats_load_library 'commands/command_writer.bash'
	bats_load_library 'commands/run/run_command_location.bash'
	write_command_requesting_input "$(run_command_location)" "$command_name"
}

assert_run_command_received_input() {
	local command_name="$1"
	local expected_input="$2"

  bats_load_library 'commands/command_writer.bash'
	bats_load_library 'commands/run/run_command_location.bash'
	assert_command_received_input "$(run_command_location)" "$command_name" "$expected_input"
}