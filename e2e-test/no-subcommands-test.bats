load 'helpers/assertDirectoryExists'
load 'helpers/assertDirectoryDoesNotExist'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/repositoriesPath'
load 'helpers/runSetup'
load 'helpers/runMrtInTestEnvironment'

setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if team json contains repositories but running 'mrt' without subcommand does not clone the repositories" {
  writeRepositoriesUrls "$(getTestingRepositoryUrl "1_TestRepository")"

  run mrt

  assert_directory_does_not_exist "$testEnvironmentDir/$(default_repositories_path)"
}