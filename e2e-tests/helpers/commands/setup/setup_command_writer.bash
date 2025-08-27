_setup_command_location() {
	bats_load_library 'fixtures/common_fixture.bash'

	echo "$(test_env_dir)/setup"
}

write_stub_setup_command() {
	local command_name="$1"
	local exit_code="$2"
	local output="$3"

	bats_load_library 'commands/command_writer.bash'
	write_stub_command "$(_setup_command_location)" "$command_name" "$exit_code" "$output"
}

write_spy_setup_command() {
	local command_name="$1"

	bats_load_library 'commands/command_writer.bash'
	write_spy_command "$(_setup_command_location)" "$command_name"
}

assert_setup_command_was_executed() {
	local command_name="$1"
	local expected_parameters="$2"

	bats_load_library 'commands/command_writer.bash'
	assert_command_was_executed_with_parameters "$(_setup_command_location)" "$command_name" "$expected_parameters"
}

assert_setup_command_was_not_executed() {
	local command_name="$1"

	bats_load_library 'commands/command_writer.bash'
	assert_command_was_not_executed "$(_setup_command_location)" "$command_name"
}

write_std_err_setup_command() {
	local command_name="$1"
	local error="$2"

	bats_load_library 'commands/command_writer.bash'
	write_std_err_command "$(_setup_command_location)" "$command_name" "$error"
}

write_setup_command_requesting_input() {
	local command_name="$1"

	bats_load_library 'commands/command_writer.bash'
	write_command_requesting_input "$(_setup_command_location)" "$command_name"
}

assert_setup_command_received_input() {
	local command_name="$1"
	local expected_input="$2"

	bats_load_library 'commands/command_writer.bash'
	assert_command_received_input "$(_setup_command_location)" "$command_name" "$expected_input"
}
