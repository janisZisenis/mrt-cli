load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeMockScript'
load 'helpers/runSetup'
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

@test "if additional setup script exists 'setup all' will execute it and pass the repository path as parameter" {
  additionalScriptsDir="$testEnvironmentDir/setup"
  setupScript="$additionalScriptsDir/setup-command/command"
  writeSpyScript "$setupScript"

  mrt setup all

  assert_spy_file_has_content "$setupScript" "$(absolutePath "$testEnvironmentDir")"
}

@test "if some additional setup script succeeds with output 'setup all' will print the script's output" {
  test_if_additional_setup_script_succeeds_setup_should_print_success_and_output "some-command"
}

@test "if another additional setup script succeeds with output 'setup all' will print the script's output" {
  test_if_additional_setup_script_succeeds_setup_should_print_success_and_output "another-command"
}

test_if_additional_setup_script_succeeds_setup_should_print_success_and_output() {
  commandName=$1
  additionalScriptsDir="$testEnvironmentDir/setup"
  setupScript="$additionalScriptsDir/$commandName/command"
  someOutput="some-output"
  writeStubScript "$setupScript" "0" "$someOutput"

  run mrt setup all

  assert_line --index 4 "Execute additional setup-script: $commandName"
  assert_line --index 5 "$someOutput"
  assert_line --index 6 "$commandName executed successfully"
  assert_line --index 7 ""
}

@test "if some additional setup script fails with output 'setup all' will print the script's output and the failure" {
  test_if_additional_setup_script_fails_setup_should_print_failure_and_output "some-command"
}

@test "if another additional setup script fails with output 'setup all' will print the script's output and the failure" {
  test_if_additional_setup_script_fails_setup_should_print_failure_and_output "another-command"
}

test_if_additional_setup_script_fails_setup_should_print_failure_and_output() {
  commandName=$1
  additionalScriptsDir="$testEnvironmentDir/setup"
  setupScript="$additionalScriptsDir/$commandName/command"
  someOutput="some-output"
  exitCode=15
  writeStubScript "$setupScript" "$exitCode" "$someOutput"

  run mrt setup all

  assert_line --index 4 "Execute additional setup-script: $commandName"
  assert_line --index 5 "$someOutput"
  assert_line --index 6 "$commandName failed with: exit status $exitCode"
  assert_line --index 7 ""
}

@test "if two additional setup scripts exist 'setup all' will execute both" {
  additionalScriptsDir="$testEnvironmentDir/setup"
  firstSetupScript="$additionalScriptsDir/setup-command1/command"
  secondSetupScript="$additionalScriptsDir/setup-command2/command"
  writeSpyScript "$firstSetupScript"
  writeSpyScript "$secondSetupScript"

  mrt setup all

  assert_spy_file_has_content "$firstSetupScript" "$(absolutePath "$testEnvironmentDir")"
  assert_spy_file_has_content "$secondSetupScript" "$(absolutePath "$testEnvironmentDir")"
}

@test "if setup script is requesting input 'setup all' should process the input" {
  additionalScriptsDir="$testEnvironmentDir/setup/input"
  additionalScriptsPath="$additionalScriptsDir/command"
  writeScriptRequestingInput "$additionalScriptsPath"
  input="some-input"

  run mrt setup all <<< $input

  assert_file_exists "$additionalScriptsDir/$input"
  assert_failure
}