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

@test "if team json does not contains any repository, 'setup' exits with error" {
  repositoriesPath=repositories
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": []
  }"

  run "$(testEnvDir)"/mrt setup

  assert_failure
  assert_output "Your team file does not contain any repositories"
}

@test "if team json does not exist, 'setup' exits with error" {
  run "$(testEnvDir)"/mrt setup

  assert_failure
  assert_output 'Could not read team file. Please make sure a "team.json" file exists next to the executable and that it follows proper JSON syntax'
}