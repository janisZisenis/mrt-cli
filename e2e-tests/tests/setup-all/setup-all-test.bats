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
	local repository_url
	repository_url="$(get_testing_repository_url "$repository")"
	write_repositories_urls "$repository_url"
	local repository_dir
	repository_dir="$(test_env_dir)/$(default_repositories_path)/$repository"
	local some_command_name="some-command"
	local another_command_name="another-command"
	write_spy_setup_command "$some_command_name"
	write_spy_setup_command "$another_command_name"

	run mrt_setup_all

	assert_line --index 0 "Start cloning repositories into \"$(default_repositories_path)\""
	assert_line --index 1 "Cloning $repository_url"
	assert_line --index 3 --regexp "Enumerating objects: [0-9]+, done."
	assert_line_reversed_output 11 "Successfully cloned $repository_url"
	assert_line_reversed_output 10 "Cloning repositories done"
	assert_line_reversed_output 9 "Installing git-hooks to repositories located in \"$(test_env_dir)/$(default_repositories_path)\""
	assert_line_reversed_output 8 "Installing git-hooks to \"$repository_dir/.git\""
	assert_line_reversed_output 7 "Done installing git-hooks to \"$repository_dir/.git\""
	assert_line_reversed_output 6 "Done installing git-hooks."
	assert_line_reversed_output 5 "Executing setup commands."
	assert_line_reversed_output 4 "Execute setup command: $another_command_name"
	assert_line_reversed_output 3 "$another_command_name executed successfully"
	assert_line_reversed_output 2 "Execute setup command: $some_command_name"
	assert_line_reversed_output 1 "$some_command_name executed successfully"
	assert_line_reversed_output 0 "Done executing setup commands."
}

@test "if setup is run without skipping git hooks it should not print skip message" {
	run mrt_setup_all

	refute_output --partial "Skipping install-git-hooks step."
}

@test "if setup command exists setup without skipping the command should not print skip message" {
	local command_name="some-command"
	write_spy_setup_command "$command_name"

	run mrt_setup_all

	refute_output --partial "Skipping setup command: $command_name"
}
