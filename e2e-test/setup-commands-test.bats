load 'helpers/common'
load 'helpers/writeMockCommand'
load 'helpers/absolutePath'
load 'helpers/executeInTestEnvironment'
load 'helpers/git'
load 'helpers/setup'


setup() {
  _common_setup
}

teardown() {
  _common_teardown
}

@test "if setup command (some-command) exists executing it will pass the team folder path as parameter" {
  test_if_setup_command_exists_executing_it_will_pass_the_team_folder_as_parameter "some-command"
}

@test "if setup command (another-command) exists executing it will pass the team folder path as parameter" {
  test_if_setup_command_exists_executing_it_will_pass_the_team_folder_as_parameter "another-command"
}

test_if_setup_command_exists_executing_it_will_pass_the_team_folder_as_parameter() {
  commandName=$1
  commandLocation="$testEnvDir/setup"
  writeSpyCommand "$commandLocation" "$commandName"

  execute setup "$commandName"

  assert_command_spy_file_has_content "$commandLocation" "$commandName" "$(absolutePath "$testEnvDir")"
}

@test "if setup command succeeds with output it will print the command's output" {
  commandName="some-command"
  someOutput="some-output"
  writeStubCommand "$testEnvDir/setup" "$commandName" "0" "$someOutput"

  run setupCommand $commandName

  assert_line --index 0 "Execute setup-command: $commandName"
  assert_line --index 1 "$someOutput"
  assert_line --index 2 "$commandName executed successfully"
}

@test "if setup command fails with output it will print the command's output and the failure" {
  commandName="another-command"
  someOutput="another-output"
  exitCode=15
  writeStubCommand "$testEnvDir/setup" "$commandName" "$exitCode" "$someOutput"

  run setupCommand "$commandName"

  assert_line --index 0 "Execute setup-command: $commandName"
  assert_line --index 1 "$someOutput"
  assert_line --index 2 "$commandName failed with: exit status $exitCode"
}

@test "if setup command is requesting input it should process the input" {
  commandName="input"
  commandLocation="$testEnvDir/setup"
  writeCommandRequestingInput "$commandLocation" "$commandName"
  input="some-input"

  run setupCommand $commandName <<< $input

  assert_command_received_input "$commandLocation" "$commandName" "$input"
}

@test "if setup command is writes to stderr it outputs stderr" {
  commandName="error"
  error="some-error"
  writeStdErrCommand "$testEnvDir/setup" "$commandName" "$error"

  run setupCommand "$commandName"

  assert_output --partial "$error"
}
