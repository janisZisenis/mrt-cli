bats_load_library 'fixtures/authenticated_fixture.bash'
bats_load_library 'git.bash'
bats_load_library 'mrt/setup.bash'
bats_load_library 'repositoriesPath.bash'
bats_load_library 'writeTeamFile.bash'

setup() {
	authenticated_setup
}

teardown() {
	authenticated_teardown
}

repositoriesPath=$(default_repositories_path)

repositoriesDir() {
	echo "$(testEnvDir)/$repositoriesPath"
}

@test "If repositories are cloned to repositories path from team file commiting on the blocked branches after setting up git-hooks should be blocked" {
	local repositoriesPath="some-path"
	local repository="1_TestRepository"
	local branchName="some-branch"
	writeRepositoriesPath "$repositoriesPath"
	writeBlockedBranches "$branchName"
	cloneTestingRepositories "$(repositoriesDir)" "$repository"
	mrtSetupGitHooks

	run commit_changes "$(repositoriesDir)/$repository" $branchName

	assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
	assert_failure
}

@test "If repositories path contains non-repository folder it does not install git-hooks" {
	local repository="1_TestRepository"
	local folderPath; folderPath="$(repositoriesDir)/$repository"
	mkdir -p "$folderPath"

	run mrtSetupGitHooks

	assert_dir_not_exist "$folderPath/.git/hooks"
}

@test "If repositories path contains 2 repositories committing on a blocked branch in the second repository after setting up git-hooks should be blocked" {
	local repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	local branchName="some-branch"
	cloneTestingRepositories "$(repositoriesDir)" "${repositories[@]}"
	writeBlockedBranches "$branchName"
	mrtSetupGitHooks

	run commit_changes "$(repositoriesDir)/${repositories[1]}" $branchName

	assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
	assert_failure
}

@test "If repositories path (some-path) does not contain any repositories setting up git-hooks prints out message that it didn't find repositories" {
	test_if_repositories_path_does_not_contain_repositories_setting_up_git_hook_prints_out_not_found_messages "some-path"
}

@test "If repositories path (another-path) does not contain any repositories setting up git-hooks prints out message that it didn't find repositories" {
	test_if_repositories_path_does_not_contain_repositories_setting_up_git_hook_prints_out_not_found_messages "some-path"
}

test_if_repositories_path_does_not_contain_repositories_setting_up_git_hook_prints_out_not_found_messages() {
	local repositoriesPath="$1"
	writeRepositoriesPath "$repositoriesPath"

	run mrtSetupGitHooks

	assert_line --index 0 "Installing git-hooks to repositories located in \"$(repositoriesDir)\""
	assert_line --index 1 "Did not find any repositories. Skip installing git-hooks."
	assert_line --index 2 "Done installing git-hooks."
}

@test "If repositories path contains two repositories setting up git-hooks prints out messages about installing the git-hooks" {
	local repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	cloneTestingRepositories "$(repositoriesDir)" "${repositories[@]}"

	run mrtSetupGitHooks

	assert_line --index 0 "Installing git-hooks to repositories located in \"$(repositoriesDir)\""
	assert_line --index 1 "Installing git-hooks to \"$(repositoriesDir)/${repositories[0]}/.git\""
	assert_line --index 2 "Done installing git-hooks to \"$(repositoriesDir)/${repositories[0]}/.git\""
	assert_line --index 3 "Installing git-hooks to \"$(repositoriesDir)/${repositories[1]}/.git\""
	assert_line --index 4 "Done installing git-hooks to \"$(repositoriesDir)/${repositories[1]}/.git\""
	assert_line --index 5 "Done installing git-hooks."
}
