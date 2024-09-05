setup() {
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'
  load 'test_helper/setupRepositories'

  _common_setup "$(testEnvDir)"
}

teardown() {
  _common_teardown "$(testEnvDir)"
}

testEnvDir() {
  echo "./testEnv"
}

@test "if team json contains existing repositories but authentication is missing, 'setup' should print message" {
  repository="1_TestRepository"

  run setupRepositories "$(testEnvDir)" "$repository"

  assert_output "You have no access to git@github-testing:janisZisenisTesting/$repository.git. Please make sure you have a valid ssh key in place."
}

@test "if team json does not contains any repository, 'setup' exits with error" {
  run setupRepositories "$(testEnvDir)" ""

  assert_failure
  assert_output "Your team file does not contain any repositories"
}

@test "if team json does not exist, 'setup' exits with error" {
  run "$(testEnvDir)"/mrt setup

  assert_failure
  assert_output 'Could not read team file. Please make sure a "team.json" file exists next to the executable and that it follows proper JSON syntax'
}