bats_load_library 'fixtures/authenticated_fixture.bash'
bats_load_library 'git.bash'
bats_load_library 'mrt/setup.bash'
bats_load_library 'repositories_path.bash'
bats_load_library 'write_team_file.bash'

setup() {
	authenticated_setup
}

teardown() {
	authenticated_teardown
}

repositories_path=$(default_repositories_path)

repositories_dir() {
	echo "$(test_env_dir)/$repositories_path"
}

@test "If repositories are cloned to repositories path from team file commiting on the blocked branches after setting up git-hooks should be blocked" {
	local repositories_path="some-path"
	local repository="1_TestRepository"
	local branch_name="some-branch"
	write_repositories_path "$repositories_path"
	write_blocked_branches "$branch_name"
	clone_testing_repositories "$(repositories_dir)" "$repository"
	mrt_setup_git_hooks

	run commit_changes "$(repositories_dir)/$repository" $branch_name

	assert_output --partial "Action \"commit\" not allowed on branch \"$branch_name\""
	assert_failure
}

@test "If repositories path contains non-repository folder it does not install git-hooks" {
	local repository="1_TestRepository"
	local folder_path; folder_path="$(repositories_dir)/$repository"
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
	clone_testing_repositories "$(repositories_dir)" "${repositories[@]}"
	write_blocked_branches "$branch_name"
	mrt_setup_git_hooks

	run commit_changes "$(repositories_dir)/${repositories[1]}" $branch_name

	assert_output --partial "Action \"commit\" not allowed on branch \"$branch_name\""
	assert_failure
}

@test "If repositories path (some-path) does not contain any repositories setting up git-hooks prints out message that it didn't find repositories" {
	test_if_repositories_path_does_not_contain_repositories_setting_up_git_hook_prints_out_not_found_messages "some-path"
}

@test "If repositories path (another-path) does not contain any repositories setting up git-hooks prints out message that it didn't find repositories" {
	test_if_repositories_path_does_not_contain_repositories_setting_up_git_hook_prints_out_not_found_messages "some-path"
}

test_if_repositories_path_does_not_contain_repositories_setting_up_git_hook_prints_out_not_found_messages() {
	local repositories_path="$1"
	write_repositories_path "$repositories_path"

	run mrt_setup_git_hooks

	assert_line --index 0 "Installing git-hooks to repositories located in \"$(repositories_dir)\""
	assert_line --index 1 "Did not find any repositories. Skip installing git-hooks."
	assert_line --index 2 "Done installing git-hooks."
}

@test "If repositories path contains two repositories setting up git-hooks prints out messages about installing the git-hooks" {
	local repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	clone_testing_repositories "$(repositories_dir)" "${repositories[@]}"

	run mrt_setup_git_hooks

	assert_line --index 0 "Installing git-hooks to repositories located in \"$(repositories_dir)\""
	assert_line --index 1 "Installing git-hooks to \"$(repositories_dir)/${repositories[0]}/.git\""
	assert_line --index 2 "Done installing git-hooks to \"$(repositories_dir)/${repositories[0]}/.git\""
	assert_line --index 3 "Installing git-hooks to \"$(repositories_dir)/${repositories[1]}/.git\""
	assert_line --index 4 "Done installing git-hooks to \"$(repositories_dir)/${repositories[1]}/.git\""
	assert_line --index 5 "Done installing git-hooks."
}
