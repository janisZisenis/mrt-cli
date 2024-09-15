load 'helpers/setup'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/repositoriesPath'
load 'helpers/runMrtInTestEnvironment'
load 'helpers/directoryAssertions'


repositoryDir() {
  echo "$testEnvironmentDir/$(default_repositories_path)/$repository"
}

setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if setup is run with skipping the clone step it should not clone the repositories" {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"
  writeRepositoriesUrls "$repositoryUrl"

  run mrt setup all --skip-clone-repositories

  assert_directory_does_not_exist "$testEnvironmentDir/$(default_repositories_path)/$repository"
}

@test "if setup is run with skipping the clone step it should print a skip message" {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"
  writeRepositoriesUrls "$repositoryUrl"

  run mrt setup all --skip-clone-repositories

  assert_line --index 0 "Skipping clone-repositories step."
}
