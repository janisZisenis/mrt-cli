load 'helpers/setupRepositories'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/commits'
load 'helpers/pushChanges'
load 'helpers/branches'

testEnvDir=$(_testEnvDir)
repository=1_TestRepository

setup() {
  _common_setup "$testEnvDir"
  authenticate

  writeRepositories "$testEnvDir" "$repository"
  "$testEnvDir"/mrt setup --skip-git-hooks
}

teardown() {
  _common_teardown "$testEnvDir"
  revoke-authentication
}

@test "After setup with 'skip-git-hooks' committing on a blocked branch is not rejected" {
  branchName="some-branch"

  writeBlockedBranches "$testEnvDir" "$branchName"

  run commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "$branchName" "some-message"

  assert_success
}

@test "After setup with 'skip-git-hooks' pushing to a blocked branch is not rejected" {
  branchName="$(unique_branch_name)"
  writeBlockedBranches "$testEnvDir" "$branchName"
  commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "$branchName" "some-message"

  run push_changes "$testEnvDir/$(default_repositories_dir)/$repository" "$branchName"

  assert_success
}

@test "After setup with 'skip-git-hooks' commiting with missing prefix in commit messages is not rejected" {
  branchName="some-branch"
  writeBlockedBranches "$testEnvDir" "$branchName"
  writeCommitPrefixRegex "$testEnvDir" "Some-Prefix"

  run commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "$branchName" "some-message"

  assert_success
}
