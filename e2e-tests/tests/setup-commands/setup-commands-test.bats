setup() {
  bats_load_library 'fixtures/common_fixture.bash'
  bats_load_library 'mrt/setup.bash'
  bats_load_library 'commands/setup/setup_command_writer.bash'

	common_setup
}

teardown() {
	common_teardown
}

@test "if setup command (some-command) exists executing it will pass the team folder path as parameter" {
	test_if_setup_command_exists_executing_it_will_pass_the_team_folder_as_parameter "some-command"
}

@test "if setup command (another-command) exists executing it will pass the team folder path as parameter" {
	test_if_setup_command_exists_executing_it_will_pass_the_team_folder_as_parameter "another-command"
}

test_if_setup_command_exists_executing_it_will_pass_the_team_folder_as_parameter() {
	local command_name="$1"
	write_spy_setup_command "$command_name"

	mrt_setup "$command_name"

	assert_setup_command_was_executed "$command_name" "$(test_env_dir)"
}

@test "if setup command succeeds with output it will print the command's output" {
	local command_name="some-command"
	local some_output="some-output"
	write_stub_setup_command "$command_name" "0" "$some_output"

	run mrt_setup $command_name

	assert_line --index 0 "Execute setup command: $command_name"
	assert_line --index 1 "$some_output"
	assert_line --index 2 "$command_name executed successfully"
}

@test "if setup command fails with output it will print the command's output and the failure" {
	local command_name="another-command"
	local some_output="another-output"
	local exit_code=15
	write_stub_setup_command "$command_name" "$exit_code" "$some_output"

	run mrt_setup "$command_name"

	assert_line --index 0 "Execute setup command: $command_name"
	assert_line --index 1 "$some_output"
	assert_line --index 2 "$command_name failed with: exit status $exit_code"
}

@test "if setup command is requesting input it should process the input" {
	local command_name="input"
	write_setup_command_requesting_input "$command_name"
	local input="some-input"

	run mrt_setup $command_name <<<$input

	assert_setup_command_received_input "$command_name" "$input"
}

@test "if setup command writes to stderr it outputs stderr" {
	local command_name="error"
	local error="some-error"
	write_std_err_setup_command "$command_name" "$error"

	run mrt_setup "$command_name"

	assert_output --partial "$error"
}
