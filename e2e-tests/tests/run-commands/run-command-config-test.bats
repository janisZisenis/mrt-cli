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
	local command_name="some-command"
	local short_description="A command that outputs some-output"
	write_dummy_run_command "$command_name"
	write_short_description "$command_name" "$short_description"

	run bats_pipe mrt_run "-h" \| grep "$command_name"

	assert_output "  $command_name $short_description"
}

@test "if command config does not contain shortDescription the default is displayed in help" {
	local command_name="some-command"
	local short_description="A command that outputs some-output"
	write_dummy_run_command "$command_name"
	write_empty_json_object_as_config "$command_name"

	run bats_pipe mrt_run "-h" \| grep "$command_name"

	assert_output "  $command_name Executes run command $command_name"
}

@test "if command config is completely empty, it should exit with an error" {
	local command_name="some-command"
	local short_description="A command that outputs some-output"
	local config_file
	config_file="$(config_file_path "$command_name")"
	write_dummy_run_command "$command_name"
	touch "$config_file"

	run mrt_run "-h"

	assert_equal "$status" 1
	assert_line --index 0 "Error while reading $config_file"
	assert_line --index 1 "While parsing config: unexpected end of JSON input"
}
