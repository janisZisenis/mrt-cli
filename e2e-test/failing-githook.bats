testEnvDir() {
  echo "./testEnv"
}

repository="1_TestRepository"

setup() {
  load 'helpers/ssh-authenticate'
  load 'helpers/common'
  load 'helpers/defaults'
  load 'helpers/setupRepositories'

  _common_setup "$(testEnvDir)"
  authenticate
  setupRepositories "$(testEnvDir)" "$repository"
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

@test "if subcommand 'git-hook' gets called with an unknown git hook it fails" {
  hookName="unknown-hook"

  run "$(testEnvDir)"/mrt git-hook --hook-name "$hookName" --repository-path "$(testEnvDir)/$(default_repositories_dir)/$repository"

  assert_output --partial "The given git-hook \"$hookName\" does not exist."
  assert_failure
}

@test "if subcommand 'git-hook' gets called with a path that does not contain a repository it fails" {
  repositoryPath="$(testEnvDir)"

  run "$(testEnvDir)"/mrt git-hook --hook-name "pre-commit" --repository-path "$repositoryPath"

  assert_output --partial "The given path \"$repositoryPath\" does not contain a repository."
  assert_failure
}

