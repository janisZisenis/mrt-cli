load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeMockScript'
load 'helpers/executeInTestEnvironment'
load 'helpers/writeTeamFile'
load 'helpers/repositoriesPath'

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

  x="[0-9]+"
  assert_line --index 0   "Start cloning repositories into \"$(default_repositories_path)\""
  assert_line --index 1   "Cloning $repositoryUrl into $(default_repositories_path)/$repository"
  assert_line --index 2 --regexp "Enumerating objects: $x, done."
  assert_line --index 3 --regexp "Counting objects: $x% \($x\/$x\), done."
  assert_line --index 4 --regexp "Compressing objects: $x% \($x\/$x\), done."
  assert_line --index 5 --regexp "Total [0-9]+ \(delta $x\), reused $x \(delta $x\), pack-reused $x \(from $x\)"
  assert_line --index 6   "Successfully cloned $repositoryUrl"
  assert_line --index 7   "Cloning repositories done"
  assert_line --index 8   "Installing git-hooks to repositories located in \"$testEnvDir/$(default_repositories_path)\""
  assert_line --index 9   "Installing git-hooks to \"$repositoryDir/.git\""
  assert_line --index 10  "Done installing git-hooks to \"$repositoryDir/.git\""
  assert_line --index 11  "Done installing git-hooks."
  assert_line --index 12  "Executing setup-scripts."
  assert_line --index 13  "Execute setup-script: $anotherScriptName"
  assert_line --index 14  "$anotherScriptName executed successfully"
  assert_line --index 15  "Execute setup-script: $someScriptName"
  assert_line --index 16  "$someScriptName executed successfully"
  assert_line --index 17  "Done executing setup-scripts."
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