load 'helpers/setup'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/repositoriesPath'
load 'helpers/git'
load 'helpers/runMrtInTestEnvironment'

repository="1_TestRepository"
repositoryUrl="$(getTestingRepositoryUrl "$repository")"
branchName="$(unique_branch_name)"

repositoryDir() {
  echo "$testEnvironmentDir/$(default_repositories_path)/$repository"
}

setup() {
  _common_setup
  authenticate

  writeRepositoriesUrls "$repositoryUrl"
  run mrt setup all --skip-install-git-hooks
}

teardown() {
  _common_teardown
  revoke-authentication
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