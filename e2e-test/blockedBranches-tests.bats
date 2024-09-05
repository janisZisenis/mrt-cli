testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'
  load 'test_helper/commits'
  load 'test_helper/pushChanges'
  load 'test_helper/defaults'
  load 'test_helper/setupRepositories'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

@test "If team json contains blocked branch, 'commiting' on the blocked branches should be blocked" {
  repository=1_TestRepository
  branchName="some-branch"
  setupRepositories "$(testEnvDir)" "$repository"
  writeBlockedBranches "$(testEnvDir)" "$branchName"

  run commit_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch, 'commiting' on another blocked branches should allowed" {
  repository=1_TestRepository
  branchName="some-branch"
  setupRepositories "$(testEnvDir)" "$repository"
  writeBlockedBranches "$(testEnvDir)" "another-branch"

  run commit_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" $branchName

  assert_success
}

@test "If team json contains 2 blocked branch, 'commiting' on second one should be blocked" {
  repository=1_TestRepository
  branchName="some-branch"
  setupRepositories "$(testEnvDir)" "$repository"
  writeBlockedBranches "$(testEnvDir)" "another-branch" "$branchName"

  run commit_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch, 'pushing' on the blocked branches should be blocked" {
  repository=1_TestRepository
  branchName="some-branch"
  setupRepositories "$(testEnvDir)" "$repository"
  writeBlockedBranches "$(testEnvDir)" "$branchName"
  commit_changes_bypassing_githooks "$(testEnvDir)/$(default_repositories_dir)/$repository" $branchName

  push_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" "$branchName"

  assert_output --partial "Action \"push\" not allowed on branch \"$branchName\""
  assert_failure
}