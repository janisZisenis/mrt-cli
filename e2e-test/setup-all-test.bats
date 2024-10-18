load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeMockScript'
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

@test "if team file contains repository and two setup scripts exist it should clone the repository, install git-hooks and execute the scripts" {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"
  writeRepositoriesUrls "$repositoryUrl"
  repositoryDir="$testEnvDir/$(default_repositories_path)/$repository"
  scriptsDir="$testEnvDir/setup"
  someScriptName="some-script"
  anotherScriptName="another-script"
  someSetupScript="$scriptsDir/$someScriptName/command"
  anotherSetupScript="$scriptsDir/$anotherScriptName/command"
  writeSpyScript "$someSetupScript"
  writeSpyScript "$anotherSetupScript"

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
  assert_line_reversed_output 5  "Executing setup-scripts."
  assert_line_reversed_output 4  "Execute setup-script: $anotherScriptName"
  assert_line_reversed_output 3  "$anotherScriptName executed successfully"
  assert_line_reversed_output 2  "Execute setup-script: $someScriptName"
  assert_line_reversed_output 1  "$someScriptName executed successfully"
  assert_line_reversed_output 0  "Done executing setup-scripts."
}

@test "if setup is run without skipping git hooks it should not print skip message" {
  run execute setup all

  refute_output --partial "Skipping install-git-hooks step."
}

@test "if setup script exists setup without skipping the script should not print skip message" {
  scriptName="some-script"
  writeSpyScript "$testEnvDir/setup/$scriptName/command"

  run execute setup all

  refute_output --partial "Skipping setup script: $scriptName"
}