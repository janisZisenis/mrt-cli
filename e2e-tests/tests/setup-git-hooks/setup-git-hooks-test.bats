bats_load_library 'ssh-authenticate'
bats_load_library 'common'
bats_load_library 'git'
bats_load_library 'setup'
bats_load_library 'repositoriesPath'
bats_load_library 'directoryAssertions'

setup() {
	common_setup
	authenticate
}

teardown() {
	revoke-authentication
	common_teardown
}

repositoriesPath=$(default_repositories_path)

repositoriesDir() {
	echo "$(testEnvDir)/$repositoriesPath"
}

@test "If repositories are cloned to repositories path from team file commiting on the blocked branches after setting up git-hooks should be blocked" {
	repositoriesPath="some-path"
	repository="1_TestRepository"
	branchName="some-branch"
	writeRepositoriesPath "$repositoriesPath"
	writeBlockedBranches "$branchName"
	cloneTestingRepositories "$(repositoriesDir)" "$repository"
	setupGitHooks

	run commit_changes "$(repositoriesDir)/$repository" $branchName

	assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
	assert_failure
}

@test "If repositories path contains non-repository folder it does not install git-hooks" {
	repository="1_TestRepository"
	folderPath="$(repositoriesDir)/$repository"
	mkdir -p "$folderPath"

	run setupGitHooks

	assert_directory_does_not_exist "$folderPath/.git/hooks"
}

@test "If repositories path contains 2 repositories committing on a blocked branch in the second repository after setting up git-hooks should be blocked" {
	repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	branchName="some-branch"
	cloneTestingRepositories "$(repositoriesDir)" "${repositories[@]}"
	writeBlockedBranches "$branchName"
	setupGitHooks

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
	repositoriesPath=$1
	writeRepositoriesPath "$repositoriesPath"

	run setupGitHooks

	assert_line --index 0 "Installing git-hooks to repositories located in \"$(repositoriesDir)\""
	assert_line --index 1 "Did not find any repositories. Skip installing git-hooks."
	assert_line --index 2 "Done installing git-hooks."
}

@test "If repositories path contains two repositories setting up git-hooks prints out messages about installing the git-hooks" {
	repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	cloneTestingRepositories "$(repositoriesDir)" "${repositories[@]}"

	run setupGitHooks

	assert_line --index 0 "Installing git-hooks to repositories located in \"$(repositoriesDir)\""
	assert_line --index 1 "Installing git-hooks to \"$(repositoriesDir)/${repositories[0]}/.git\""
	assert_line --index 2 "Done installing git-hooks to \"$(repositoriesDir)/${repositories[0]}/.git\""
	assert_line --index 3 "Installing git-hooks to \"$(repositoriesDir)/${repositories[1]}/.git\""
	assert_line --index 4 "Done installing git-hooks to \"$(repositoriesDir)/${repositories[1]}/.git\""
	assert_line --index 5 "Done installing git-hooks."
}
