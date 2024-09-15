load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/commits'
load 'helpers/pushChanges'
load 'helpers/repositoriesPath'
load 'helpers/runSetup'
load 'helpers/branches'
load 'helpers/clone'
load 'helpers/assertDirectoryDoesNotExist'

setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "If repositories path contains 2 repositories committing on a blocked branch in the second repository after setting up git-hooks should be blocked" {
  repositories=(
    "1_TestRepository"
    "2_TestRepository"
  )
  cloneTestingRepositories "$testEnvironmentDir/$(default_repositories_path)" "${repositories[@]}"
  branchName="some-branch"
  writeBlockedBranches "$branchName"
  setupGitHooks

  run commit_changes "$testEnvironmentDir/$(default_repositories_path)/${repositories[1]}" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}