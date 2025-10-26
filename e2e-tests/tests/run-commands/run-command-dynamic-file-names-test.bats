setup() {
	bats_load_library 'fixtures/common_fixture.bash'
	bats_load_library 'mrt/run.bash'
  bats_load_library 'commands/command_writer.bash'
	bats_load_library 'commands/run/run_command_location.bash'
  bats_load_library "scripts/script_writer.bash"
  bats_load_library "commands/run/run_command_config_writer.bash"
  bats_load_library "commands/run/run_command_writer.bash"
  bats_load_library "scripts/script_writer.bash"

	common_setup
}

teardown() {
	common_teardown
}

@test "if some-command is run it should execute it" {
	local command_name="some-command"
  local parameters=("some" "--flag")
  local command_file_name="command.bash"
  local command_location="$(run_command_location)"
  write_spy_script "$command_location/$command_name/$command_file_name"
	write_command_file_name "$command_name" "$command_file_name"

  run mrt_run "$command_name" -- "${parameters[@]}"

	assert_script_was_executed_with_parameters "$(test_env_dir)/run/$command_name/$command_file_name" "$(test_env_dir) ${parameters[*]}"
}