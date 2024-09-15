load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/commits'
load 'helpers/setup'
load 'helpers/clone'
load 'helpers/repositoriesPath'


setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

repositoriesDir() {
  echo "$testEnvironmentDir/$(default_repositories_path)"
}

@test "If team json contains blocked branch commiting on the blocked branches after setting up git-hooks should be blocked" {
  repository="1_TestRepository"
  branchName="some-branch"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  writeBlockedBranches "$branchName"
  setupGitHooks

  run commit_changes "$(repositoriesDir)/$repository" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch commiting on another blocked branch after setting up git-hooks should be allowed" {
  repository="1_TestRepository"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  writeBlockedBranches "some-branch"
  setupGitHooks

  run commit_changes "$(repositoriesDir)/$repository" "another-branch"

  assert_success
}

@test "If team json contains 2 blocked branch commiting on second one after setting up git-hooks should be blocked" {
  repository="1_TestRepository"
  branchName="some-branch"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  writeBlockedBranches "another-branch" "$branchName"
  setupGitHooks

  run commit_changes "$(repositoryDir)" $branchName

  assert_output --partial "Action 1 \"commit\" not allowed on branch \"$branchName\""
  assert_failure
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