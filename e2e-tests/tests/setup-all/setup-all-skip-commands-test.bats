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
	local someCommandName="some-command"
	local anotherCommandName="another-command"
	write_spy_setup_command "$someCommandName"
	write_spy_setup_command "$anotherCommandName"

	run mrt_setup_all "--skip-$someCommandName"

	assert_setup_command_was_not_executed "$someCommandName"
	assert_setup_command_was_executed "$anotherCommandName" "$(test_env_dir)"
}

@test "if two setup commands exist setup all with skipping the second it should only run the first" {
	local someCommandName="some-command"
	local anotherCommandName="another-command"
	write_spy_setup_command "$someCommandName"
	write_spy_setup_command "$anotherCommandName"

	run mrt_setup_all "--skip-$anotherCommandName"

	assert_setup_command_was_executed "$someCommandName" "$(test_env_dir)"
	assert_setup_command_was_not_executed "$anotherCommandName"
}

@test "if one setup commands exists setup all with skipping the command prints out skip message" {
	local commandName="some-command"
	write_spy_setup_command "$commandName"

	run mrt_setup_all "--skip-$commandName"

	assert_output --partial "Skipping setup command: $commandName"
}
