load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeMockCommand'
load 'helpers/executeInTestEnvironment'
load 'helpers/writeTeamFile'
load 'helpers/repositoriesPath'
load 'helpers/assertLineReversed'

setup() {
  _common_setup
  authenticate
}

teardown() {
  revoke-authentication
  _common_teardown
}

@test "if team file contains repository and two setup commands exist it should clone the repository, install git-hooks and execute the commands" {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"
  writeRepositoriesUrls "$repositoryUrl"
  repositoryDir="$testEnvDir/$(default_repositories_path)/$repository"
  someCommandName="some-command"
  anotherCommandName="another-command"
  commandLocation="$testEnvDir/setup"
  writeSpyCommand "$commandLocation" "$someCommandName"
  writeSpyCommand "$commandLocation" "$anotherCommandName"

  run execute setup all

  assert_line --index 0   "Start cloning repositories into \"$(default_repositories_path)\""
  assert_line --index 1   "Cloning $repositoryUrl into $(default_repositories_path)/$repository"
  assert_line --index 2 --regexp "Enumerating objects: [0-9]+, done."
  assert_line_reversed_output 11 "Successfully cloned $repositoryUrl"
  assert_line_reversed_output 10 "Cloning repositories done"
  assert_line_reversed_output 9  "Installing git-hooks to repositories located in \"$testEnvDir/$(default_repositories_path)\""
  assert_line_reversed_output 8  "Installing git-hooks to \"$repositoryDir/.git\""
  assert_line_reversed_output 7  "Done installing git-hooks to \"$repositoryDir/.git\""
  assert_line_reversed_output 6  "Done installing git-hooks."
  assert_line_reversed_output 5  "Executing setup-commands."
  assert_line_reversed_output 4  "Execute setup-command: $anotherCommandName"
  assert_line_reversed_output 3  "$anotherCommandName executed successfully"
  assert_line_reversed_output 2  "Execute setup-command: $someCommandName"
  assert_line_reversed_output 1  "$someCommandName executed successfully"
  assert_line_reversed_output 0  "Done executing setup-commands."
}

@test "if setup is run without skipping git hooks it should not print skip message" {
  run execute setup all

  refute_output --partial "Skipping install-git-hooks step."
}

@test "if setup command exists setup without skipping the command should not print skip message" {
  commandName="some-command"
  writeSpyCommand "$testEnvDir/setup" "$commandName"

  run execute setup all

  refute_output --partial "Skipping setup command: $commandName"
}