load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/repositoriesPath'
load 'helpers/setup'
load 'helpers/runMrtInTestEnvironment'

repository="1_TestRepository"

setup() {
  _common_setup
  authenticate
  setupAll "$repository"
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if subcommand 'git-hook' gets called with an unknown git hook it fails" {
  hookName="unknown-hook"
  repositoryDir="$testEnvironmentDir/$(default_repositories_path)/$repository"

  run mrt git-hook --hook-name "$hookName" --repository-path "$repositoryDir"

  assert_output --partial "The given git-hook \"$hookName\" does not exist."
  assert_failure
}

@test "if subcommand 'git-hook' gets called with a path that does not contain a repository it fails" {
  repositoryPath="$testEnvironmentDir"

  run mrt git-hook --hook-name "pre-commit" --repository-path "$repositoryPath"

  assert_output --partial "The given path \"$repositoryPath\" does not contain a repository."
  assert_failure
}

