testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'test_helper/assertDirectoryExists'
  load 'test_helper/assertDirectoryDoesNotExist'
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'
  load 'test_helper/commitChanges'
  load 'test_helper/pushChanges'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

defaultRepositoriesPath="repositories"

@test "If team json contains blocked branch, 'commiting' on the blocked branches should be blocked" {
  repository=1_TestRepository
  branchName="some-branch"
  writeTeamFile "$(testEnvDir)" "{
      \"blockedBranches\": [
        \"$branchName\"
      ],
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch, 'commiting' on another blocked branches should allowed" {
  repository=1_TestRepository
  branchName="some-branch"
  writeTeamFile "$(testEnvDir)" "{
      \"blockedBranches\": [
        \"another-branch\"
      ],
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" $branchName

  assert_success
}

@test "If team json contains 2 blocked branch, 'commiting' on second one should be blocked" {
  repository=1_TestRepository
  branchName="some-branch"
  writeTeamFile "$(testEnvDir)" "{
      \"blockedBranches\": [
        \"another-branch\",
        \"$branchName\"
      ],
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" $branchName

  assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}

@test "If team json contains blocked branch, 'pushing' on the blocked branches should be blocked" {
  repository=1_TestRepository
  branchName="some-branch"
  writeTeamFile "$(testEnvDir)" "{
      \"blockedBranches\": [
        \"$branchName\"
      ],
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup
  commit_changes_bypassing_githooks "$(testEnvDir)/$defaultRepositoriesPath/$repository" $branchName

  push_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" $branchName

  assert_output --partial "Action \"push\" not allowed on branch \"$branchName\""
  assert_failure
}