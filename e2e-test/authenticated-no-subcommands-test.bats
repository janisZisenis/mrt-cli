load 'helpers/assertDirectoryExists'
load 'helpers/assertDirectoryDoesNotExist'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/setupRepositories'
load 'helpers/runMrtInTestEnvironment'

repositoriesDir() {
  echo "$testEnvironmentDir/$(default_repositories_dir)"
}

setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if team json contains repositories but running without 'setup all' does not clone the repositories" {
  writeRepositories "$(getTestingRepositoryUrl "1_TestRepository")"

  run mrt

  assert_directory_does_not_exist "$(repositoriesDir)"
}