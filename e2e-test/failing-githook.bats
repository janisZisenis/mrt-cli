load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/setupRepositories'
load 'helpers/runMrtInTestEnvironment'

repository="1_TestRepository"
repositoryUrl="$(getTestingRepositoryUrl "$repository")"

repositoriesDir(){
  echo "$testEnvironmentDir/$(default_repositories_dir)/$repository"
}

setup() {
  _common_setup
  authenticate
  setupRepositories "$repositoryUrl"
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if subcommand 'git-hook' gets called with an unknown git hook it fails" {
  hookName="unknown-hook"

  run mrt git-hook --hook-name "$hookName" --repository-path "$(repositoriesDir)"

  assert_output --partial "The given git-hook \"$hookName\" does not exist."
  assert_failure
}

@test "if subcommand 'git-hook' gets called with a path that does not contain a repository it fails" {
  repositoryPath="$testEnvironmentDir"

  run mrt git-hook --hook-name "pre-commit" --repository-path "$repositoryPath"

  assert_output --partial "The given path \"$repositoryPath\" does not contain a repository."
  assert_failure
}

