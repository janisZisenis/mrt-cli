setup() {
	bats_load_library 'mrt/setup.bash'
	bats_load_library 'fixtures/common_fixture.bash'
	bats_load_library 'commands/setup/setup_command_writer.bash'

	common_setup
}

teardown() {
	common_teardown
}

@test "if two setup commands exist setup all with skipping the first it should only run the second" {
	local some_command_name="some-command"
	local another_command_name="another-command"
	write_spy_setup_command "$some_command_name"
	write_spy_setup_command "$another_command_name"

	run mrt_setup_all "--skip-$some_command_name"

	assert_setup_command_was_not_executed "$some_command_name"
	assert_setup_command_was_executed "$another_command_name" "$(test_env_dir)"
}

@test "if two setup commands exist setup all with skipping the second it should only run the first" {
	local some_command_name="some-command"
	local another_command_name="another-command"
	write_spy_setup_command "$some_command_name"
	write_spy_setup_command "$another_command_name"

	run mrt_setup_all "--skip-$another_command_name"

	assert_setup_command_was_executed "$some_command_name" "$(test_env_dir)"
	assert_setup_command_was_not_executed "$another_command_name"
}

@test "if one setup commands exists setup all with skipping the command prints out skip message" {
	local command_name="some-command"
	write_spy_setup_command "$command_name"

	run mrt_setup_all "--skip-$command_name"

	assert_output --partial "Skipping setup command: $command_name"
}
