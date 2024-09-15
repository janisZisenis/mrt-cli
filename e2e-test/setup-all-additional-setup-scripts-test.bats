load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeMockScript'
load 'helpers/setup'
load 'helpers/absolutePath'
load 'helpers/runMrtInTestEnvironment'
load 'helpers/assertFileExists'

setup() {
  _common_setup
  authenticate

  setupAll "1_TestRepository"
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if two additional setup scripts exist it will execute both" {
  additionalScriptsDir="$testEnvironmentDir/setup"
  someSetupScript="$additionalScriptsDir/some-command/command"
  anotherSetupScript="$additionalScriptsDir/another-command/command"
  writeSpyScript "$someSetupScript"
  writeSpyScript "$anotherSetupScript"

  mrt setup all

  assert_spy_file_has_content "$someSetupScript" "$(absolutePath "$testEnvironmentDir")"
  assert_spy_file_has_content "$anotherSetupScript" "$(absolutePath "$testEnvironmentDir")"
}
