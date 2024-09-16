load 'helpers/common'
load 'helpers/writeMockScript'
load 'helpers/run'
load 'helpers/absolutePath'


setup() {
  _common_setup
}

teardown() {
  _common_teardown
}

@test "if some-script is run it should execute it" {
  test_if_run_is_executed_with_script_name_it_should_execute_the_specified_script "some-script"
}

@test "if another-script is run it should execute it" {
  test_if_run_is_executed_with_script_name_it_should_execute_the_specified_script "another-script"
}

test_if_run_is_executed_with_script_name_it_should_execute_the_specified_script() {
  scriptName=$1
  scriptPath="$testEnvironmentDir/run/$scriptName/command"
  writeSpyScript "$scriptPath"

  run runScript "$scriptName"

  assert_spy_file_has_content "$scriptPath" "$(absolutePath "$testEnvironmentDir")"
}

@test "if script succeeds with output it will print the script's output" {
  scriptName="some-script"
  someOutput="some-output"
  writeStubScript "$testEnvironmentDir/run/$scriptName/command" "0" "$someOutput"

  run runScript "$scriptName"

  assert_output "$someOutput"
}

@test "if script is requesting input it should process the input" {
  scriptName="input"
  additionalScriptsDir="$testEnvironmentDir/run/$scriptName"
  additionalScriptsPath="$additionalScriptsDir/command"
  writeScriptRequestingInput "$additionalScriptsPath"
  input="some-input"

  run runScript $scriptName <<< $input

  assert_file_exists "$additionalScriptsDir/$input"
}

@test "if script writes to stderr it outputs stderr" {
  scriptName="error"
  error="some-error"
  writeStdErrScript "$testEnvironmentDir/run/$scriptName/command" "$error"

  run runScript "$scriptName"

  assert_output "$error"
}