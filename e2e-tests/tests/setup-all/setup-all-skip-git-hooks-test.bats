set_fixture_variables() {
  bats_load_library 'git.bash'
  bats_load_library 'repositories_path.bash'
  bats_load_library 'fixtures/common_fixture.bash'

  repository="1_TestRepository"
  branch_name="$(unique_branch_name)"
  repository_dir="$(test_env_dir)/$(default_repositories_path)/$repository"
}

setup() {
  bats_load_library 'mrt/setup.bash'
  bats_load_library 'write_team_file.bash'
  bats_load_library 'test_repositories.bash'
  bats_load_library 'fixtures/authenticated_fixture.bash'

  set_fixture_variables
  authenticated_setup

	write_repositories_urls "$(get_testing_repository_url "$repository")"
	mrt_setup_all --skip-install-git-hooks
}

teardown() {
	authenticated_teardown
}

@test "After setup all with 'skip-git-hooks' committing on a blocked branch is not rejected" {
	write_blocked_branches "$branch_name"

	run commit_changes "$repository_dir" "$branch_name" "some-message"

	assert_success
}

@test "After setup all with 'skip-git-hooks' pushing to a blocked branch is not rejected" {
	write_blocked_branches "$branch_name"
	commit_changes "$repository_dir" "$branch_name" "some-message"

	run push_changes "$repository_dir" "$branch_name"

	assert_success
}

@test "After setup all with 'skip-git-hooks' commiting with missing prefix in commit messages is not rejected" {
	write_blocked_branches "$branch_name"
	write_commit_prefix_regex "Some-Prefix"

	run commit_changes "$repository_dir" "$branch_name" "some-message"

	assert_success
}
