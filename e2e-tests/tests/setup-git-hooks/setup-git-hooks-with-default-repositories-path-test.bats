set_fixture_variables() {
  bats_load_library 'fixtures/common_fixture.bash'
  bats_load_library 'repositories_path.bash'

  repositories_path=$(default_repositories_path)
  repositories_dir="$(test_env_dir)/$repositories_path"
}

setup() {
  bats_load_library 'fixtures/authenticated_fixture.bash'
  bats_load_library 'git.bash'
  bats_load_library 'mrt/setup.bash'
  bats_load_library 'repositories_path.bash'
  bats_load_library 'write_team_file.bash'

  set_fixture_variables
	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "If repositories path contains non-repository folder it does not install git-hooks" {
	local repository="1_TestRepository"
	local folder_path; folder_path="$repositories_dir/$repository"
	mkdir -p "$folder_path"

	run mrt_setup_git_hooks

	assert_dir_not_exist "$folder_path/.git/hooks"
}

@test "If repositories path contains 2 repositories committing on a blocked branch in the second repository after setting up git-hooks should be blocked" {
	local repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	local branch_name="some-branch"
	clone_testing_repositories "$repositories_dir" "${repositories[@]}"
	write_blocked_branches "$branch_name"
	mrt_setup_git_hooks

	run commit_changes "$repositories_dir/${repositories[1]}" $branch_name

	assert_output --partial "Action \"commit\" not allowed on branch \"$branch_name\""
	assert_failure
}

@test "If repositories path contains two repositories setting up git-hooks prints out messages about installing the git-hooks" {
	local repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	clone_testing_repositories "$repositories_dir" "${repositories[@]}"

	run mrt_setup_git_hooks

	assert_line --index 0 "Installing git-hooks to repositories located in \"$repositories_dir\""
	assert_line --index 1 "Installing git-hooks to \"$repositories_dir/${repositories[0]}/.git\""
	assert_line --index 2 "Done installing git-hooks to \"$repositories_dir/${repositories[0]}/.git\""
	assert_line --index 3 "Installing git-hooks to \"$repositories_dir/${repositories[1]}/.git\""
	assert_line --index 4 "Done installing git-hooks to \"$repositories_dir/${repositories[1]}/.git\""
	assert_line --index 5 "Done installing git-hooks."
}
