bats_load_library 'setup'
bats_load_library 'ssh-authenticate'
bats_load_library 'common_fixture'
bats_load_library 'repositoriesPath'
bats_load_library 'git'
bats_load_library 'testRepositories'

repository="1_TestRepository"
repositoryUrl="$(getTestingRepositoryUrl "$repository")"
branchName="$(unique_branch_name)"

repositoryDir() {
	echo "$(testEnvDir)/$(default_repositories_path)/$repository"
}

setup() {
	common_setup
	authenticate

	writeRepositoriesUrls "$repositoryUrl"
	run execute setup all --skip-install-git-hooks
}

teardown() {
	revoke-authentication
	common_teardown
}

@test "After setup all with 'skip-git-hooks' committing on a blocked branch is not rejected" {
	writeBlockedBranches "$branchName"

	run commit_changes "$(repositoryDir)" "$branchName" "some-message"

	assert_success
}

@test "After setup all with 'skip-git-hooks' pushing to a blocked branch is not rejected" {
	writeBlockedBranches "$branchName"
	commit_changes "$(repositoryDir)" "$branchName" "some-message"

	run push_changes "$(repositoryDir)" "$branchName"

	assert_success
}

@test "After setup all with 'skip-git-hooks' commiting with missing prefix in commit messages is not rejected" {
	writeBlockedBranches "$branchName"
	writeCommitPrefixRegex "Some-Prefix"

	run commit_changes "$(repositoryDir)" "$branchName" "some-message"

	assert_success
}

@test "If setup is run with skipping git hooks, it should print skip message" {
	assert_line --partial "Skipping install-git-hooks step."
}
