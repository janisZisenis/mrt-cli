bats_load_library 'setup'
bats_load_library 'fixtures/common_fixture'
bats_load_library 'repositoriesPath'
bats_load_library 'commands/setup/setupCommandWriter'

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

	run execute setup all "--skip-$someCommandName"

	assert_setup_command_was_not_executed "$someCommandName"
	assert_setup_command_was_executed "$anotherCommandName" "$(testEnvDir)"
}

@test "if two setup commands exist setup all with skipping the second it should only run the first" {
	someCommandName="some-command"
	anotherCommandName="another-command"
	writeSpySetupCommand "$someCommandName"
	writeSpySetupCommand "$anotherCommandName"

	run execute setup all "--skip-$anotherCommandName"

	assert_setup_command_was_executed "$someCommandName" "$(testEnvDir)"
	assert_setup_command_was_not_executed "$anotherCommandName"
}

@test "if one setup commands exists setup all with skipping the command prints out skip message" {
	commandName="some-command"
	writeSpySetupCommand "$commandName"

	run execute setup all "--skip-$commandName"

	assert_output --partial "Skipping setup command: $commandName"
}
