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
  scriptsDir="$testEnvironmentDir/run/$scriptName"
  scriptsPath="$scriptsDir/command"
  writeScriptRequestingInput "$scriptsPath"
  input="some-input"

  run runScript $scriptName <<< $input

  assert_file_exists "$scriptsDir/$input"
}

@test "if script writes to stderr it outputs stderr" {
  scriptName="error"
  error="some-error"
  writeStdErrScript "$testEnvironmentDir/run/$scriptName/command" "$error"

  run runScript "$scriptName"

  assert_output "$error"
}

@test "if script fails with code 1 it will fail with error code 1 as well" {
  scriptName="some-script"
  exitCode=1
  writeStubScript "$testEnvironmentDir/run/$scriptName/command" "$exitCode" ""

  run runScript "$scriptName"

  assert_equal "$status" "$exitCode"
  assert_failure
}

@test "if script fails with code 2 it will fail with error code 2 as well" {
  scriptName="some-script"
  exitCode=2
  writeStubScript "$testEnvironmentDir/run/$scriptName/command" "$exitCode" ""

  run runScript "$scriptName"

  assert_equal "$status" "$exitCode"
  assert_failure
}

@test "if script exits with code 0 it will succeed" {
  scriptName="some-script"
  exitCode=0
  writeStubScript "$testEnvironmentDir/run/$scriptName/command" "$exitCode" ""

  run runScript "$scriptName"

  assert_success
}