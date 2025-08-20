bats_load_library 'common'
bats_load_library 'ssh-authenticate'
bats_load_library 'writeTeamFile'
bats_load_library 'repositoriesPath'
bats_load_library 'assertLineReversed'
bats_load_library 'commands/setupCommandWriter'
bats_load_library 'testRepositories'

setup() {
	common_setup
	authenticate
}

teardown() {
	revoke-authentication
	common_teardown
}

@test "if team file contains repository and two setup commands exist it should clone the repository, install git-hooks and execute the commands" {
	repository="1_TestRepository"
	repositoryUrl="$(getTestingRepositoryUrl "$repository")"
	writeRepositoriesUrls "$repositoryUrl"
	repositoryDir="$(testEnvDir)/$(default_repositories_path)/$repository"
	someCommandName="some-command"
	anotherCommandName="another-command"
	writeSpySetupCommand "$someCommandName"
	writeSpySetupCommand "$anotherCommandName"

	run execute setup all

	assert_line --index 0 "Start cloning repositories into \"$(default_repositories_path)\""
	assert_line --index 1 "Cloning $repositoryUrl"
	assert_line --index 3 --regexp "Enumerating objects: [0-9]+, done."
	assert_line_reversed_output 11 "Successfully cloned $repositoryUrl"
	assert_line_reversed_output 10 "Cloning repositories done"
	assert_line_reversed_output 9 "Installing git-hooks to repositories located in \"$(testEnvDir)/$(default_repositories_path)\""
	assert_line_reversed_output 8 "Installing git-hooks to \"$repositoryDir/.git\""
	assert_line_reversed_output 7 "Done installing git-hooks to \"$repositoryDir/.git\""
	assert_line_reversed_output 6 "Done installing git-hooks."
	assert_line_reversed_output 5 "Executing setup commands."
	assert_line_reversed_output 4 "Execute setup command: $anotherCommandName"
	assert_line_reversed_output 3 "$anotherCommandName executed successfully"
	assert_line_reversed_output 2 "Execute setup command: $someCommandName"
	assert_line_reversed_output 1 "$someCommandName executed successfully"
	assert_line_reversed_output 0 "Done executing setup commands."
}

@test "if setup is run without skipping git hooks it should not print skip message" {
	run execute setup all

	refute_output --partial "Skipping install-git-hooks step."
}

@test "if setup command exists setup without skipping the command should not print skip message" {
	commandName="some-command"
	writeSpyCommand "$(testEnvDir)/setup" "$commandName"

	run execute setup all

	refute_output --partial "Skipping setup command: $commandName"
}
