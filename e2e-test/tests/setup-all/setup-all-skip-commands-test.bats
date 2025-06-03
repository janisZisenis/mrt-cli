load '../../helpers/setup'
load '../../helpers/ssh-authenticate'
load '../../helpers/common'
load '../../helpers/repositoriesPath'
load '../../helpers/directoryAssertions'
load '../../helpers/writeMockCommand'
load '../../helpers/absolutePath'

setup() {
	common_setup
}

teardown() {
	common_teardown
}

@test "if two setup commands exist setup all with skipping the first it should only run the second" {
	someCommandName="some-command"
	anotherCommandName="another-command"
	commandLocation="$(testEnvDir)/setup"
	writeSpyCommand "$commandLocation" "$someCommandName"
	writeSpyCommand "$commandLocation" "$anotherCommandName"

	run execute setup all "--skip-$someCommandName"

	assert_command_spy_file_does_not_exist "$commandLocation" "$someCommandName"
	assert_command_spy_file_exists "$commandLocation" "$anotherCommandName"
}

@test "if two setup commands exist setup all with skipping the second it should only run the first" {
	someCommandName="some-command"
	anotherCommandName="another-command"
	commandLocation="$(testEnvDir)/setup"
	writeSpyCommand "$commandLocation" "$someCommandName"
	writeSpyCommand "$commandLocation" "$anotherCommandName"

	run execute setup all "--skip-$anotherCommandName"

	assert_command_spy_file_exists "$commandLocation" "$someCommandName"
	assert_command_spy_file_does_not_exist "$commandLocation" "$anotherCommandName"
}

@test "if one setup commands exists setup all with skipping the command prints out skip message" {
	commandName="some-command"
	commandLocation="$(testEnvDir)/setup"
	writeSpyCommand "$commandLocation" "$commandName"

	run execute setup all "--skip-$commandName"

	assert_output --partial "Skipping setup command: $commandName"
}
