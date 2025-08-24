setup() {
  bats_load_library 'write_team_file.bash'
  bats_load_library 'repositories_path.bash'
  bats_load_library 'test_repositories.bash'
  bats_load_library 'assert_line_reversed.bash'
  bats_load_library 'commands/setup/setup_command_writer.bash'
  bats_load_library 'fixtures/authenticated_fixture.bash'
  bats_load_library 'mrt/setup.bash'

	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "if team file contains repository and two setup commands exist it should clone the repository, install git-hooks and execute the commands" {
	local repository="1_TestRepository"
	local repositoryUrl; repositoryUrl="$(get_testing_repository_url "$repository")"
	write_repositories_urls "$repositoryUrl"
	local repositoryDir; repositoryDir="$(test_env_dir)/$(default_repositories_path)/$repository"
	local someCommandName="some-command"
	local anotherCommandName="another-command"
	write_spy_setup_command "$someCommandName"
	write_spy_setup_command "$anotherCommandName"

	run mrt_setup_all

	assert_line --index 0 "Start cloning repositories into \"$(default_repositories_path)\""
	assert_line --index 1 "Cloning $repositoryUrl"
	assert_line --index 3 --regexp "Enumerating objects: [0-9]+, done."
	assert_line_reversed_output 11 "Successfully cloned $repositoryUrl"
	assert_line_reversed_output 10 "Cloning repositories done"
	assert_line_reversed_output 9 "Installing git-hooks to repositories located in \"$(test_env_dir)/$(default_repositories_path)\""
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
	run mrt_setup_all

	refute_output --partial "Skipping install-git-hooks step."
}

@test "if setup command exists setup without skipping the command should not print skip message" {
	local commandName="some-command"
	write_spy_setup_command "$commandName"

	run mrt_setup_all

	refute_output --partial "Skipping setup command: $commandName"
}
