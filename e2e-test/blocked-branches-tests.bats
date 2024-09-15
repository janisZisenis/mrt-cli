load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/commits'
load 'helpers/pushChanges'
load 'helpers/repositoriesPath'
load 'helpers/runSetup'
load 'helpers/branches'
load 'helpers/clone'
load 'helpers/assertDirectoryDoesNotExist'

repository="1_TestRepository"
repositoriesPath=$(default_repositories_path)

repositoriesDir() {
  echo "$testEnvironmentDir/$repositoriesPath"
}

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

@test "If team json contains blocked branch commiting on the blocked branches after setting up git-hooks should be blocked" {
  branchName="some-branch"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  writeBlockedBranches "$branchName"
  setupGitHooks

  run commit_changes "$(repositoryDir)" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If repositories are cloned to specific path commiting on the blocked branches after setting up git-hooks should be blocked" {
  branchName="some-branch"
  repositoriesPath="some-path"
  writeRepositoriesPath "$repositoriesPath"
  writeBlockedBranches "$branchName"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  setupGitHooks

  run commit_changes "$(repositoryDir)" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch commiting on another blocked branch after setting up git-hooks should be allowed" {
  branchName="some-branch"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  writeBlockedBranches "another-branch"
  setupGitHooks

  run commit_changes "$(repositoryDir)" $branchName

  assert_success
}

@test "If team json contains 2 blocked branch commiting on second one after setting up git-hooks should be blocked" {
  branchName="some-branch"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  writeBlockedBranches "another-branch" "$branchName"
  setupGitHooks

  run commit_changes "$(repositoryDir)" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch pushing on the blocked after setting up git-hooks branch should be blocked" {
  branchName="$(unique_branch_name)"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  writeBlockedBranches "$branchName"
  commit_changes_bypassing_githooks "$(repositoryDir)" "$branchName"
  setupGitHooks

  push_changes "$(repositoryDir)" "$branchName"

  assert_output --partial "Action \"push\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If repositories path contains non-repository folder it does not install githooks" {
  folderPath="$(repositoriesDir)/$repository"
  mkdir -p "$folderPath"

  run setupGitHooks

  assert_directory_does_not_exist "$folderPath/.git/hooks"
}