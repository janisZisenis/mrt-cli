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

teardown() {
  _common_teardown "$testEnvDir"
  revoke-authentication
}

@test "if additional setup script exists 'setup' will execute it and pass the repository path as parameter" {
  additionalScriptsDir="$testEnvDir/setup"
  setupScript="$additionalScriptsDir/setup-command/command"
  writeSpyScript "$setupScript"

  "$testEnvDir"/mrt setup

  assert_spy_file_has_content "$setupScript" "$(absolutePath "$testEnvDir")"
}

@test "if some additional setup script succeeds with output 'setup' will print the script's output" {
  test_if_additional_setup_script_succeeds_setup_should_print_success_and_output "some-command"
}

@test "if another additional setup script succeeds with output 'setup' will print the script's output" {
  test_if_additional_setup_script_succeeds_setup_should_print_success_and_output "another-command"
}

test_if_additional_setup_script_succeeds_setup_should_print_success_and_output() {
  commandName=$1
  additionalScriptsDir="$testEnvDir/setup"
  setupScript="$additionalScriptsDir/$commandName/command"
  someOutput="some-output"
  writeStubScript "$setupScript" "0" "$someOutput"

  run "$testEnvDir"/mrt setup

  assert_line --index 1 "Execute additional setup-script: $commandName"
  assert_line --index 2 "$someOutput"
  assert_line --index 3 "$commandName executed successfully"
  assert_line --index 4 ""
}

@test "if some additional setup script fails with output 'setup' will print the script's output and the failure" {
  test_if_additional_setup_script_fails_setup_should_print_failure_and_output "some-command"
}

@test "if another additional setup script fails with output 'setup' will print the script's output and the failure" {
  test_if_additional_setup_script_fails_setup_should_print_failure_and_output "another-command"
}

test_if_additional_setup_script_fails_setup_should_print_failure_and_output() {
  commandName=$1
  additionalScriptsDir="$testEnvDir/setup"
  setupScript="$additionalScriptsDir/$commandName/command"
  someOutput="some-output"
  exitCode=15
  writeStubScript "$setupScript" "$exitCode" "$someOutput"

  run "$testEnvDir"/mrt setup

  assert_line --index 1 "Execute additional setup-script: $commandName"
  assert_line --index 2 "$someOutput"
  assert_line --index 3 "$commandName failed with: exit status $exitCode"
  assert_line --index 4 ""
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
