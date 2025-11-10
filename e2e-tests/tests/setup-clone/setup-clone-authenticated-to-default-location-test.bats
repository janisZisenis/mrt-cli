repositories_dir() {
	bats_load_library 'fixtures/common_fixture.bash'
	bats_load_library 'repositories_path.bash'

	echo "$(test_env_dir)/$(default_repositories_path)"
}

cloned_git_folder() {
	local repository="$1"

	echo "$(repositories_dir)/$repository/.git"
}

setup() {
	bats_load_library 'fixtures/authenticated_fixture.bash'
	bats_load_library 'test_repositories.bash'
	bats_load_library 'mrt/clone.bash'
	bats_load_library 'git.bash'
	bats_load_library 'assert_line_reversed.bash'
	bats_load_library 'write_team_file.bash'

	authenticated_setup
}

teardown() {
	authenticated_teardown
}
