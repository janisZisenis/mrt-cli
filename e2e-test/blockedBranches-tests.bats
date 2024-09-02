testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'test_helper/assertDirectoryExists'
  load 'test_helper/assertDirectoryDoesNotExist'
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

defaultRepositoriesPath="repositories"

@test "If team json contains 2 blocked branches, 'commiting' on the first blocked branches should be blocked" {
  repository=1_TestRepository
  branchName="some-feature-branch"
  writeTeamFile "$(testEnvDir)" "{
      \"blockedBranches\": [
        \"$branchName\",
        \"another-feature-branch\"
      ],
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" checkout -b $branchName
  touch "$(testEnvDir)/$defaultRepositoriesPath/$repository"/some_file
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" add .

  run git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" commit -m "Should be rejected"

  assert_output "Action \"commit\" not allowed on branch \"$branchName\""
  assert_failure
}