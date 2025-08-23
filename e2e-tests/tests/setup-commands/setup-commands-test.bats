setup() {
  bats_load_library 'fixtures/common_fixture.bash'
  bats_load_library 'mrt/setup.bash'
  bats_load_library 'commands/setup/setupCommandWriter.bash'

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
	writeSpySetupCommand "$commandName"

	mrtSetup "$commandName"

	assert_setup_command_was_executed "$commandName" "$(testEnvDir)"
}

@test "if setup command succeeds with output it will print the command's output" {
	commandName="some-command"
	someOutput="some-output"
	writeStubSetupCommand "$commandName" "0" "$someOutput"

	run mrtSetup $commandName

	assert_line --index 0 "Execute setup command: $commandName"
	assert_line --index 1 "$someOutput"
	assert_line --index 2 "$commandName executed successfully"
}

@test "if setup command fails with output it will print the command's output and the failure" {
	commandName="another-command"
	someOutput="another-output"
	exitCode=15
	writeStubSetupCommand "$commandName" "$exitCode" "$someOutput"

	run mrtSetup "$commandName"

	assert_line --index 0 "Execute setup command: $commandName"
	assert_line --index 1 "$someOutput"
	assert_line --index 2 "$commandName failed with: exit status $exitCode"
}

@test "if setup command is requesting input it should process the input" {
	commandName="input"
	writeSetupCommandRequestingInput "$commandName"
	input="some-input"

	run mrtSetup $commandName <<<$input

	assert_setup_command_received_input "$commandName" "$input"
}

@test "if setup command writes to stderr it outputs stderr" {
	commandName="error"
	error="some-error"
	writeStdErrSetupCommand "$commandName" "$error"

	run mrtSetup "$commandName"

	assert_output --partial "$error"
}
