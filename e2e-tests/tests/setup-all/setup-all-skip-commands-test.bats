bats_load_library 'mrt/clone.bash'
bats_load_library 'mrt/setup.bash'
bats_load_library 'fixtures/common_fixture.bash'
bats_load_library 'repositoriesPath.bash'
bats_load_library 'commands/setup/setupCommandWriter.bash'

setup() {
	common_setup
}

teardown() {
	common_teardown
}

@test "if two setup commands exist setup all with skipping the first it should only run the second" {
	someCommandName="some-command"
	anotherCommandName="another-command"
	writeSpySetupCommand "$someCommandName"
	writeSpySetupCommand "$anotherCommandName"

	run mrtSetupAll "--skip-$someCommandName"

	assert_setup_command_was_not_executed "$someCommandName"
	assert_setup_command_was_executed "$anotherCommandName" "$(testEnvDir)"
}

@test "if two setup commands exist setup all with skipping the second it should only run the first" {
	someCommandName="some-command"
	anotherCommandName="another-command"
	writeSpySetupCommand "$someCommandName"
	writeSpySetupCommand "$anotherCommandName"

	run mrtSetupAll "--skip-$anotherCommandName"

	assert_setup_command_was_executed "$someCommandName" "$(testEnvDir)"
	assert_setup_command_was_not_executed "$anotherCommandName"
}

@test "if one setup commands exists setup all with skipping the command prints out skip message" {
	commandName="some-command"
	writeSpySetupCommand "$commandName"

	run mrtSetupAll "--skip-$commandName"

	assert_output --partial "Skipping setup command: $commandName"
}
