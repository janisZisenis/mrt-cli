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

@test "--skip-clone-repositories" {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"
  writeRepositoriesUrls "$repositoryUrl"

  run mrt setup all --skip-clone-repositories

  assert_directory_does_not_exist "$testEnvironmentDir/$(default_repositories_path)/$repository"
}

