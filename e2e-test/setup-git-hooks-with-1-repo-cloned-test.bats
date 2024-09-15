load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/commits'
load 'helpers/setup'
load 'helpers/clone'
load 'helpers/repositoriesPath'
load 'helpers/branch'
load 'helpers/push'

repository="1_TestRepository"

setup() {
  _common_setup
  authenticate

  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  setupGitHooks
}

teardown() {
  _common_teardown
  revoke-authentication
}

repositoriesDir() {
  echo "$testEnvironmentDir/$(default_repositories_path)"
}

@test "If team json contains blocked branch commiting on the blocked branches after setting up git-hooks should be blocked" {
  branchName="some-branch"
  writeBlockedBranches "$branchName"

  run commit_changes "$(repositoriesDir)/$repository" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch commiting on another blocked branch after setting up git-hooks should be allowed" {
  writeBlockedBranches "some-branch"

  run commit_changes "$(repositoriesDir)/$repository" "another-branch"

  assert_success
}

@test "If team json contains 2 blocked branch commiting on second one after setting up git-hooks should be blocked" {
  branchName="some-branch"
  writeBlockedBranches "another-branch" "$branchName"

  run commit_changes "$(repositoriesDir)/$repository" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch pushing on the blocked after setting up git-hooks branch should be blocked" {
  branchName="$(unique_branch_name)"
  commit_changes "$(repositoriesDir)/$repository" "$branchName"
  writeBlockedBranches "$branchName"

  push_changes "$(repositoriesDir)/$repository" "$branchName"

  assert_output --partial "Action \"push\" not allowed on branch \"$branchName\""
  assert_failure
}
