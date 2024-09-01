setup() {
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'

  _common_setup "$(testEnvDir)"
}

teardown() {
  _common_teardown "$(testEnvDir)"
}

testEnvDir() {
  echo "./testEnv"
}

@test "if team json contains existing repositories but authentication is missing, 'setup' should print message" {
  repositoriesPath=repositories
  repositoryUrl=git@github-testing:janisZisenisTesting/1_TestRepository.git
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"$repositoryUrl\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_output "You have no access to $repositoryUrl. Please make sure you have a valid ssh key in place."
}
