load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/commits'
load 'helpers/pushChanges'
load 'helpers/repositoriesPath'
load 'helpers/setupRepositories'
load 'helpers/branches'

repository="1_TestRepository"
repositoryUrl="$(getTestingRepositoryUrl "$repository")"
repositoriesPath=$(default_repositories_path)

repositoryDir() {
  echo "$testEnvironmentDir/$repositoriesPath/$repository"
}

setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "If team json contains blocked branch, 'commiting' on the blocked branches should be blocked" {
  branchName="some-branch"
  setupRepositories "$repositoryUrl"
  writeBlockedBranches "$branchName"

  run commit_changes "$(repositoryDir)" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If repositories are cloned to specific repositories path and team json contains blocked branch, 'commiting' on the blocked branches should be blocked" {
  branchName="some-branch"
  repositoriesPath="some-path"
  writeRepositoriesPath "$repositoriesPath"
  setupRepositories "$repositoryUrl"
  writeBlockedBranches "$branchName"

  run commit_changes "$(repositoryDir)" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch, 'commiting' on another blocked branches should allowed" {
  branchName="some-branch"
  setupRepositories "$repositoryUrl"
  writeBlockedBranches "another-branch"

  run commit_changes "$(repositoryDir)" $branchName

  assert_success
}

@test "If team json contains 2 blocked branch, 'commiting' on second one should be blocked" {
  branchName="some-branch"
  setupRepositories "$repositoryUrl"
  writeBlockedBranches "another-branch" "$branchName"

  run commit_changes "$(repositoryDir)" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch, 'pushing' on the blocked branches should be blocked" {
  branchName="$(unique_branch_name)"
  setupRepositories "$repositoryUrl"
  writeBlockedBranches "$branchName"
  commit_changes_bypassing_githooks "$(repositoryDir)" "$branchName"

  push_changes "$(repositoryDir)" "$branchName"

  assert_output --partial "Action \"push\" not allowed on branch \"$branchName\""
  assert_failure
}