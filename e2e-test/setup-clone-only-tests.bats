load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeTeamFile'
load 'helpers/runMrtInTestEnvironment'
load 'helpers/commits'
load 'helpers/assertDirectoryExists'
load 'helpers/repositoriesPath'
load 'helpers/runSetup'
load 'helpers/branches'
load 'helpers/pushChanges'

repository="1_TestRepository"
repositoryUrl=$(getTestingRepositoryUrl "$repository")

repositoryDir() {
  echo "$testEnvironmentDir/$(default_repositories_path)/$repository"
}

setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if team json does not contain repositoriesPath 'setup all' clones repository into 'repositories' folder" {
  writeRepositoriesUrls "$repositoryUrl"

  run mrt setup clone-repositories

  assert_directory_exists "$(repositoryDir)/.git"
}

@test "After setup with 'clone-repositories' subcommand, committing on a blocked branch is not rejected" {
  writeRepositoriesUrls "$repositoryUrl"
  mrt setup clone-repositories
  branchName="some-branch"
  writeBlockedBranches "$branchName"

  run commit_changes "$(repositoryDir)" "$branchName" "some-message"

  assert_success
}

@test "After setup with 'clone-repositories' subcommand, pushing to a blocked branch is not rejected" {
  writeRepositoriesUrls "$repositoryUrl"
  mrt setup clone-repositories
  branchName="$(unique_branch_name)-some-branch"
  writeBlockedBranches "$branchName"
  commit_changes "$(repositoryDir)" "$branchName" "some-message"

  run push_changes "$(repositoryDir)" "$branchName"

  assert_success
}

@test "After setup with 'clone-repositories' subcommand, commiting with missing prefix in commit messages is not rejected" {
  writeRepositoriesUrls "$repositoryUrl"
  mrt setup clone-repositories
  branchName="$(unique_branch_name)-some-branch"
  writeBlockedBranches "$branchName"
  writeCommitPrefixRegex "Some-Prefix"

  run commit_changes "$(repositoryDir)" "$branchName" "some-message"

  assert_success
}