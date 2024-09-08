load 'helpers/assertDirectoryExists'
load 'helpers/assertDirectoryDoesNotExist'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/setupRepositories'

testEnvDir=$(_testEnvDir)

setup() {
  _common_setup "$testEnvDir"
  authenticate
}

teardown() {
  _common_teardown "$testEnvDir"
  revoke-authentication
}

@test "if team json contains repositories but running without 'setup all' does not clone the repositories" {
  writeRepositories "$testEnvDir" "$(getTestingRepositoryUrl "1_TestRepository")"

  run "$testEnvDir"/mrt

  assert_directory_does_not_exist "$testEnvDir/$(default_repositories_dir)"
}