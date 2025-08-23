bats_load_library 'fixtures/common_fixture'
bats_load_library 'commands/run/runCommandWriter'
bats_load_library 'mrt/run'

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
	local commandName=$1; shift
	local parameters=("$@")
	writeSpyRunCommand "$commandName"

	run mrtRun "$commandName" -- "${parameters[@]}"

	assert_run_command_was_executed_with_parameters "$commandName" "$(testEnvDir) ${parameters[*]}"
}

@test "if command succeeds with output it will print the command's output" {
	commandName="some-command"
	someOutput="some-output"
	writeStubRunCommand "$commandName" "0" "$someOutput"

	run mrtRun "$commandName"

	assert_output "$someOutput"
}

@test "if command is requesting input it should process the input" {
	commandName="input"
	writeRunCommandRequestingInput "$commandName"
	input="some-input"

	run mrtRun $commandName <<<$input

	assert_run_command_received_input "$commandName" "$input"
}

@test "if command writes to stderr it outputs stderr" {
	commandName="error"
	error="some-error"
	writeStdErrRunCommand "$commandName" "$error"

	run mrtRun "$commandName"

	assert_output "$error"
}

# shellcheck disable=SC2030
@test "if command fails with code 1 it will fail with error code 1 as well" {
  test_if_command_fails_with_error_code_it_fails_with_the_same_error_code 1
}

# shellcheck disable=SC2030
@test "if command fails with code 2 it will fail with error code 2 as well" {
  test_if_command_fails_with_error_code_it_fails_with_the_same_error_code 2
}

test_if_command_fails_with_error_code_it_fails_with_the_same_error_code() {
  local exitCode="$1"
  local commandName="some-command"
  writeStubRunCommand "$commandName" "$exitCode" ""

  run mrtRun "$commandName"

  # shellcheck disable=SC2031
  assert_equal "$status" "$exitCode"
  assert_failure
}

@test "if command exits with code 0 it will succeed" {
	commandName="some-command"
	exitCode=0
	writeStubRunCommand "$commandName" "$exitCode" ""

	run mrtRun "$commandName"

	assert_success
}
