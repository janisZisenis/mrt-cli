setup() {
  bats_load_library 'fixtures/common_fixture.bash'
  bats_load_library 'commands/run/run_command_writer.bash'
  bats_load_library 'mrt/run.bash'

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
	write_spy_run_command "$commandName"

	run mrt_run "$commandName" -- "${parameters[@]}"

	assert_run_command_was_executed_with_parameters "$commandName" "$(test_env_dir) ${parameters[*]}"
}

@test "if command succeeds with output it will print the command's output" {
	local commandName="some-command"
	local someOutput="some-output"
	write_stub_run_command "$commandName" "0" "$someOutput"

	run mrt_run "$commandName"

	assert_output "$someOutput"
}

@test "if command is requesting input it should process the input" {
	local commandName="input"
	write_run_command_requesting_input "$commandName"
	local input="some-input"

	run mrt_run $commandName <<<$input

	assert_run_command_received_input "$commandName" "$input"
}

@test "if command writes to stderr it outputs stderr" {
	local commandName="error"
	local error="some-error"
	write_std_err_run_command "$commandName" "$error"

	run mrt_run "$commandName"

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
  write_stub_run_command "$commandName" "$exitCode" ""

  run mrt_run "$commandName"

  # shellcheck disable=SC2031
  assert_equal "$status" "$exitCode"
  assert_failure
}

@test "if command exits with code 0 it will succeed" {
	local commandName="some-command"
	local exitCode=0
	write_stub_run_command "$commandName" "$exitCode" ""

	run mrt_run "$commandName"

	assert_success
}
