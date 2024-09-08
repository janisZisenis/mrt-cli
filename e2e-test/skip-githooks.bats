load 'helpers/setupRepositories'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/commits'
load 'helpers/pushChanges'
load 'helpers/branches'
load 'helpers/runMrtInTestEnvironment'

repository="1_TestRepository"
repositoryUrl="$(getTestingRepositoryUrl "$repository")"
branchName="$(unique_branch_name)"

repositoryDir() {
  echo "$testEnvironmentDir/$(default_repositories_dir)/$repository"
}

setup() {
  _common_setup
  authenticate

  writeRepositories "$repositoryUrl"
  mrt setup all --skip-git-hooks
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "After setup with 'skip-git-hooks' committing on a blocked branch is not rejected" {
  writeBlockedBranches "$branchName"

  run commit_changes "$(repositoryDir)" "$branchName" "some-message"

  assert_success
}

@test "After setup with 'skip-git-hooks' pushing to a blocked branch is not rejected" {
  writeBlockedBranches "$branchName"
  commit_changes "$(repositoryDir)" "$branchName" "some-message"

  run push_changes "$(repositoryDir)" "$branchName"

  assert_success
}

@test "After setup with 'skip-git-hooks' commiting with missing prefix in commit messages is not rejected" {
  writeBlockedBranches "$branchName"
  writeCommitPrefixRegex "Some-Prefix"

  run commit_changes "$(repositoryDir)" "$branchName" "some-message"

  assert_success
}
