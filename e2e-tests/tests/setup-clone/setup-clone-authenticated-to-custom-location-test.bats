setup() {
  bats_load_library 'fixtures/authenticated_fixture.bash'
  bats_load_library 'fixtures/common_fixture.bash'
  bats_load_library 'mrt/clone.bash'
  bats_load_library 'write_team_file.bash'

	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "if team json contains repositories_path it clones the repositories into given repositories_path folder" {
	local repositories_path="xyz"
	write_repositories_path "$repositories_path"
	local repository="1_TestRepository"

	run clone_repositories_using_mrt "$repository"

	assert_dir_exist "$(test_env_dir)/$repositories_path/$repository/.git"
}
