bats_load_library 'common'
bats_load_library 'commandWriter'
bats_load_library 'git'
bats_load_library 'setup'
bats_load_library 'setupCommandWriter'

setup() {
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
	commandName=$1
	writeSpySetupCommand "$commandName"

	execute setup "$commandName"

	assert_setup_command_was_executed "$commandName" "$(testEnvDir)"
}

@test "if setup command succeeds with output it will print the command's output" {
	commandName="some-command"
	someOutput="some-output"
	writeStubSetupCommand "$commandName" "0" "$someOutput"

	run setupCommand $commandName

	assert_line --index 0 "Execute setup command: $commandName"
	assert_line --index 1 "$someOutput"
	assert_line --index 2 "$commandName executed successfully"
}

@test "if setup command fails with output it will print the command's output and the failure" {
	commandName="another-command"
	someOutput="another-output"
	exitCode=15
	writeStubSetupCommand "$commandName" "$exitCode" "$someOutput"

	run setupCommand "$commandName"

	assert_line --index 0 "Execute setup command: $commandName"
	assert_line --index 1 "$someOutput"
	assert_line --index 2 "$commandName failed with: exit status $exitCode"
}

@test "if setup command is requesting input it should process the input" {
	commandName="input"
	commandLocation="$(testEnvDir)/setup"
	writeCommandRequestingInput "$commandLocation" "$commandName"
	input="some-input"

	run setupCommand $commandName <<<$input

	assert_command_received_input "$commandLocation" "$commandName" "$input"
}

@test "if setup command writes to stderr it outputs stderr" {
	commandName="error"
	error="some-error"
	writeStdErrSetupCommand "$commandName" "$error"

	run setupCommand "$commandName"

	assert_output --partial "$error"
}
