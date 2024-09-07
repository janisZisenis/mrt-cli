load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeMockScript'
load 'helpers/setupRepositories'
load 'helpers/absolutePath'

testEnvDir="$(_testEnvDir)"
repository="1_TestRepository"

setup() {
  _common_setup "$testEnvDir"
  authenticate

  setupRepositories "$testEnvDir" "$repository"
}

#teardown() {
#  _common_teardown "$testEnvDir"
#  revoke-authentication
#}

@test "if additional setup script exists 'setup' will execute it and pass the repository path as parameter" {
  additionalScriptsDir="$testEnvDir/setup"
  setupScript="$additionalScriptsDir/setup-command/command"
  writeSpyScript "$setupScript"

  "$testEnvDir"/mrt setup

  assert_spy_file_has_content "$setupScript" "$(absolutePath "$testEnvDir")"
}

@test "if additional setup script has output 'setup' will have the same output" {
  additionalScriptsDir="$testEnvDir/setup"
  setupScript="$additionalScriptsDir/setup-command/command"
  someOutput="some-output"
  writeStubScript "$setupScript" "0" "$someOutput"

  run "$testEnvDir"/mrt setup

  assert_output --partial "$someOutput"
}

@test "if two additional setup scripts exist 'setup' will execute both" {
  additionalScriptsDir="$testEnvDir/setup"
  firstSetupScript="$additionalScriptsDir/setup-command1/command"
  secondSetupScript="$additionalScriptsDir/setup-command2/command"
  writeSpyScript "$firstSetupScript"
  writeSpyScript "$secondSetupScript"

  "$testEnvDir"/mrt setup

  assert_spy_file_has_content "$firstSetupScript" "$(absolutePath "$testEnvDir")"
  assert_spy_file_has_content "$secondSetupScript" "$(absolutePath "$testEnvDir")"
}

