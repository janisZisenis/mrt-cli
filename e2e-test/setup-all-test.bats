load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeMockScript'
load 'helpers/runMrtInTestEnvironment'
load 'helpers/writeTeamFile'
load 'helpers/repositoriesPath'
load 'helpers/absolutePath'

setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if team file contains repository and two additional setup scripts exist it should clone the repository, install git-hooks and execute the scripts" {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"
  writeRepositoriesUrls "$repositoryUrl"
  repositoryDir="$testEnvironmentDir/$(default_repositories_path)/$repository"
  additionalScriptsDir="$testEnvironmentDir/setup"
  someScriptName="some-script"
  anotherScriptName="another-script"
  someSetupScript="$additionalScriptsDir/$someScriptName/command"
  anotherSetupScript="$additionalScriptsDir/$anotherScriptName/command"
  writeSpyScript "$someSetupScript"
  writeSpyScript "$anotherSetupScript"

  run mrt setup all

  x="[0-9]+"
  assert_line --index 0   "Start cloning repositories into \"$(default_repositories_path)\""
  assert_line --index 1   "Cloning $repositoryUrl into $(default_repositories_path)/$repository"
  assert_line --index 2 --regexp "Enumerating objects: $x, done."
  assert_line --index 3 --regexp "Counting objects: $x% \($x\/$x\), done."
  assert_line --index 4 --regexp "Compressing objects: $x% \($x\/$x\), done."
  assert_line --index 5 --regexp "Total [0-9]+ \(delta $x\), reused $x \(delta $x\), pack-reused $x \(from $x\)"
  assert_line --index 6   "Successfully cloned $repositoryUrl"
  assert_line --index 7   "Cloning repositories done"
  assert_line --index 8   "Installing git-hooks to repositories located in \"$(absolutePath "$testEnvironmentDir/$(default_repositories_path)")\""
  assert_line --index 9   "Installing git-hooks to \"$(absolutePath "$repositoryDir")/.git\""
  assert_line --index 10  "Done installing git-hooks to \"$(absolutePath "$repositoryDir")/.git\""
  assert_line --index 11  "Done installing git-hooks."
  assert_line --index 12  "Executing additional setup-scripts."
  assert_line --index 13  "Execute additional setup-script: $anotherScriptName"
  assert_line --index 14  "$anotherScriptName executed successfully"
  assert_line --index 15  "Execute additional setup-script: $someScriptName"
  assert_line --index 16  "$someScriptName executed successfully"
  assert_line --index 17  "Done executing additional setup-scripts."
}
