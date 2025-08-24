setup() {
  bats_load_library 'fixtures/common_fixture.bash'
  bats_load_library 'commands/run/run_command_writer.bash'
  bats_load_library 'commands/run/run_command_config_writer.bash'
  bats_load_library 'mrt/run.bash'

	common_setup
}

teardown() {
	common_teardown
}

@test "if command config contains shortDescription, it is displayed in help" {
	local commandName="some-command"
  local shortDescription="A command that outputs some-output"
  write_dummy_run_command "$commandName"
  write_short_description "$commandName" "$shortDescription"

  run bats_pipe mrt_run "-h" \| grep "$commandName"

	assert_output "  $commandName $shortDescription"
}

@test "if command config does not contain shortDescription the default is displayed in help" {
	local commandName="some-command"
  local shortDescription="A command that outputs some-output"
  write_dummy_run_command "$commandName"
  write_empty_json_object_as_config "$commandName"

  run bats_pipe mrt_run "-h" \| grep "$commandName"

	assert_output "  $commandName Executes run command $commandName"
}

@test "if command config is completely empty, it should exit with an error" {
	local commandName="some-command"
  local shortDescription="A command that outputs some-output"
  local configFile; configFile="$(config_file_path "$commandName")"
  write_dummy_run_command "$commandName"
  touch "$configFile"

  run mrt_run "-h"

  assert_equal "$status" 1
  assert_line --index 0 "Error while reading $configFile"
  assert_line --index 1 "While parsing config: unexpected end of JSON input"
}