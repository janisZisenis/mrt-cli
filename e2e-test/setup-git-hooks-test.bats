load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/git'
load 'helpers/setup'
load 'helpers/repositoriesPath'
load 'helpers/assertDirectoryDoesNotExist'


setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

repositoriesPath=$(default_repositories_path)

repositoriesDir() {
  echo "$testEnvironmentDir/$repositoriesPath"
}

@test "If repositories are cloned to repositories path from team file commiting on the blocked branches after setting up git-hooks should be blocked" {
  repositoriesPath="some-path"
  repository="1_TestRepository"
  branchName="some-branch"
  writeRepositoriesPath "$repositoriesPath"
  writeBlockedBranches "$branchName"
  cloneTestingRepositories "$(repositoriesDir)" "$repository"
  setupGitHooks

  run commit_changes "$(repositoriesDir)/$repository" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If repositories path contains non-repository folder it does not install git-hooks" {
  repository="1_TestRepository"
  folderPath="$(repositoriesDir)/$repository"
  mkdir -p "$folderPath"

  run setupGitHooks

  assert_directory_does_not_exist "$folderPath/.git/hooks"
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

