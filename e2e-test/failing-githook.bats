testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

defaultRepositoriesPath="repositories"

@test "if subcommand 'git-hook' gets called with an unknown git hook it fails" {
  repository=1_TestRepository
  hookName="unknown-hook"
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run "$(testEnvDir)"/mrt git-hook --hook-name "$hookName" --repository-path "$(testEnvDir)/$defaultRepositoriesPath/$repository"

  assert_output --partial "The given git-hook \"$hookName\" does not exist."
  assert_failure
}

@test "if subcommand 'git-hook' gets called with a path that does not contain a repository it fails" {
  repository=1_TestRepository
  repositoryPath="$(testEnvDir)"
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run "$(testEnvDir)"/mrt git-hook --hook-name "pre-commit" --repository-path "$repositoryPath"

  assert_output --partial "The given path \"$repositoryPath\" does not contain a repository."
  assert_failure
}

