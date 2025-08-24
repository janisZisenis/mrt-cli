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
	local commandName="$1"
	write_spy_setup_command "$commandName"

	mrt_setup "$commandName"

	assert_setup_command_was_executed "$commandName" "$(test_env_dir)"
}

@test "if setup command succeeds with output it will print the command's output" {
	local commandName="some-command"
	local someOutput="some-output"
	write_stub_setup_command "$commandName" "0" "$someOutput"

	run mrt_setup $commandName

	assert_line --index 0 "Execute setup command: $commandName"
	assert_line --index 1 "$someOutput"
	assert_line --index 2 "$commandName executed successfully"
}

@test "if setup command fails with output it will print the command's output and the failure" {
	local commandName="another-command"
	local someOutput="another-output"
	local exitCode=15
	write_stub_setup_command "$commandName" "$exitCode" "$someOutput"

	run mrt_setup "$commandName"

	assert_line --index 0 "Execute setup command: $commandName"
	assert_line --index 1 "$someOutput"
	assert_line --index 2 "$commandName failed with: exit status $exitCode"
}

@test "if setup command is requesting input it should process the input" {
	local commandName="input"
	write_setup_command_requesting_input "$commandName"
	local input="some-input"

	run mrt_setup $commandName <<<$input

	assert_setup_command_received_input "$commandName" "$input"
}

@test "if setup command writes to stderr it outputs stderr" {
	local commandName="error"
	local error="some-error"
	write_std_err_setup_command "$commandName" "$error"

	run mrt_setup "$commandName"

	assert_output --partial "$error"
}
