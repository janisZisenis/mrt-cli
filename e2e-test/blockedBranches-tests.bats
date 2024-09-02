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
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" checkout -b $branchName
  touch "$(testEnvDir)/$defaultRepositoriesPath/$repository"/some_file
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" add .

  run git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" commit -m "Should be rejected"

  assert_output "Action \"commit\" not allowed on branch \"$branchName\""
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
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" checkout -b $branchName
  touch "$(testEnvDir)/$defaultRepositoriesPath/$repository"/some_file
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" add .

  run git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" commit -m "Should be rejected"

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
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" checkout -b $branchName
  touch "$(testEnvDir)/$defaultRepositoriesPath/$repository"/some_file
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" add .

  run git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" commit -m "Should be rejected"

  assert_output "Action \"commit\" not allowed on branch \"$branchName\""
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
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" checkout -b $branchName
  touch "$(testEnvDir)/$defaultRepositoriesPath/$repository"/some_file
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" add .
  git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" commit -m "Some Commit" --no-verify

  run git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" push --set-upstream origin "$branchName"
  if [[ $(git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" ls-remote --heads origin "$branchName") ]]
  then
    git -C "$(testEnvDir)/$defaultRepositoriesPath/$repository" push origin --delete "$branchName"
  fi

  assert_output --partial "Action \"push\" not allowed on branch \"$branchName\""
  assert_failure
}