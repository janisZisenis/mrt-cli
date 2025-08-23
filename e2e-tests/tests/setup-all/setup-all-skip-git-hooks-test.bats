bats_load_library 'mrt/setup.bash'
bats_load_library 'repositoriesPath.bash'
bats_load_library 'writeTeamFile.bash'
bats_load_library 'git.bash'
bats_load_library 'testRepositories.bash'
bats_load_library 'fixtures/common_fixture.bash'
bats_load_library 'fixtures/authenticated_fixture.bash'

set_fixture_variables() {
  repository="1_TestRepository"
  branchName="$(unique_branch_name)"
  repositoryDir="$(testEnvDir)/$(default_repositories_path)/$repository"
}

setup() {
  set_fixture_variables
  authenticated_setup

	writeRepositoriesUrls "$(getTestingRepositoryUrl "$repository")"
	mrtSetupAll --skip-install-git-hooks
}

teardown() {
	authenticated_teardown
}

@test "After setup all with 'skip-git-hooks' committing on a blocked branch is not rejected" {
	writeBlockedBranches "$branchName"

	run commit_changes "$repositoryDir" "$branchName" "some-message"

	assert_success
}

@test "After setup all with 'skip-git-hooks' pushing to a blocked branch is not rejected" {
	writeBlockedBranches "$branchName"
	commit_changes "$repositoryDir" "$branchName" "some-message"

	run push_changes "$repositoryDir" "$branchName"

	assert_success
}

@test "After setup all with 'skip-git-hooks' commiting with missing prefix in commit messages is not rejected" {
	writeBlockedBranches "$branchName"
	writeCommitPrefixRegex "Some-Prefix"

	run commit_changes "$repositoryDir" "$branchName" "some-message"

	assert_success
}
