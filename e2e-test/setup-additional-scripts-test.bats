load 'helpers/common'
load 'helpers/writeMockScript'
load 'helpers/absolutePath'
load 'helpers/runMrtInTestEnvironment'
load 'helpers/git'
load 'helpers/setup'


setup() {
  _common_setup
}

teardown() {
  _common_teardown
}

@test "if additional setup script (some-script) exists executing it will pass the team folder path as parameter" {
  test_if_additional_setup_script_exists_executing_it_will_pass_the_team_folder_as_parameter "some-script"
}

@test "if additional setup script (another-script) exists executing it will pass the team folder path as parameter" {
  test_if_additional_setup_script_exists_executing_it_will_pass_the_team_folder_as_parameter "another-script"
}

test_if_additional_setup_script_exists_executing_it_will_pass_the_team_folder_as_parameter() {
  scriptName=$1
  additionalScriptsDir="$testEnvironmentDir/setup"
  scriptPath="$additionalScriptsDir/$scriptName/command"
  writeSpyScript "$scriptPath"

  mrt setup "$scriptName"

  assert_spy_file_has_content "$scriptPath" "$(absolutePath "$testEnvironmentDir")"
}

@test "if additional setup script succeeds with output it will print the script's output" {
  scriptName="some-script"
  someOutput="some-output"
  writeStubScript "$testEnvironmentDir/setup/$scriptName/command" "0" "$someOutput"

  run setupScript $scriptName

  assert_line --index 0 "Execute additional setup-script: $scriptName"
  assert_line --index 1 "$someOutput"
  assert_line --index 2 "$scriptName executed successfully"
}

@test "if additional setup script fails with output it will print the script's output and the failure" {
  scriptName="another-script"
  someOutput="another-output"
  exitCode=15
  writeStubScript "$testEnvironmentDir/setup/$scriptName/command" "$exitCode" "$someOutput"

  run setupScript "$scriptName"

  assert_line --index 0 "Execute additional setup-script: $scriptName"
  assert_line --index 1 "$someOutput"
  assert_line --index 2 "$scriptName failed with: exit status $exitCode"
}

@test "if setup script is requesting input it should process the input" {
  scriptName="input"
  additionalScriptsDir="$testEnvironmentDir/setup/$scriptName"
  additionalScriptsPath="$additionalScriptsDir/command"
  writeScriptRequestingInput "$additionalScriptsPath"
  input="some-input"

  run setupScript $scriptName <<< $input

  assert_file_exists "$additionalScriptsDir/$input"
}

@test "if setup script is writes to stderr it outputs stderr" {
  scriptName="error"
  error="some-error"
  writeStdErrScript "$testEnvironmentDir/setup/$scriptName/command" "$error"

  run setupScript "$scriptName"

  assert_output --partial "$error"
}
