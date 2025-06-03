load '../../helpers/common'
load '../../helpers/writeMockCommand'
load '../../helpers/absolutePath'

setup() {
	common_setup
}

teardown() {
	common_teardown
}

@test "if some-command is run it should execute it" {
	test_if_run_is_executed_with_command_name_it_should_pass_root_dir_and_parameters_to_it "some-command" "some" "--flag"
}

@test "if another-command is run it should execute it" {
	test_if_run_is_executed_with_command_name_it_should_pass_root_dir_and_parameters_to_it "another-command" "another" "parameter"
}

test_if_run_is_executed_with_command_name_it_should_pass_root_dir_and_parameters_to_it() {
	commandName=$1
	shift
	parameters=("$@")
	commandLocation="$(testEnvDir)/run"
	writeSpyCommand "$(testEnvDir)/run" "$commandName"

	run runCommand "$commandName" -- "${parameters[@]}"

	assert_command_spy_file_has_content "$commandLocation" "$commandName" "$(absolutePath "$(testEnvDir)") ${parameters[*]}"
}

@test "if command succeeds with output it will print the command's output" {
	commandName="some-command"
	someOutput="some-output"
	commandLocation="$(testEnvDir)/run"
	writeStubCommand "$commandLocation" "$commandName" "0" "$someOutput"

	run runCommand "$commandName"

	assert_output "$someOutput"
}

@test "if command is requesting input it should process the input" {
	commandName="input"
	commandLocation="$(testEnvDir)/run"
	writeCommandRequestingInput "$commandLocation" "$commandName"
	input="some-input"

	run runCommand $commandName <<<$input

	assert_command_received_input "$commandLocation" "$commandName" "$input"
}

@test "if command writes to stderr it outputs stderr" {
	commandName="error"
	error="some-error"
	commandLocation="$(testEnvDir)/run"
	writeStdErrCommand "$commandLocation" "$commandName" "$error"

	run runCommand "$commandName"

	assert_output "$error"
}

@test "if command fails with code 1 it will fail with error code 1 as well" {
	commandName="some-command"
	exitCode=1
	writeStubCommand "$(testEnvDir)/run" "$commandName" "$exitCode" ""

	run runCommand "$commandName"

	assert_equal "$status" "$exitCode"
	assert_failure
}

@test "if command fails with code 2 it will fail with error code 2 as well" {
	commandName="some-command"
	exitCode=2
	writeStubCommand "$(testEnvDir)/run" "$commandName" "$exitCode" ""

	run runCommand "$commandName"

	assert_equal "$status" "$exitCode"
	assert_failure
}

@test "if command exits with code 0 it will succeed" {
	commandName="some-command"
	exitCode=0
	writeStubCommand "$(testEnvDir)/run" "$commandName" "$exitCode" ""

	run runCommand "$commandName"

	assert_success
}
