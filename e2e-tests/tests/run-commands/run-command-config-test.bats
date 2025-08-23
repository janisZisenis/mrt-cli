setup() {
  bats_load_library 'fixtures/common_fixture.bash'
  bats_load_library 'commands/run/runCommandWriter.bash'
  bats_load_library 'commands/run/runCommandConfigWriter.bash'
  bats_load_library 'mrt/run.bash'

	common_setup
}

teardown() {
	common_teardown
}

@test "if command config contains shortDescription, it is displayed in help" {
	local commandName="some-command"
  local shortDescription="A command that outputs some-output"
  writeDummyRunCommand "$commandName"
  writeShortDescription "$commandName" "$shortDescription"

  run bats_pipe mrtRun "-h" \| grep "$commandName"

	assert_output "  $commandName $shortDescription"
}

@test "if command config does not contain shortDescription the default is displayed in help" {
	local commandName="some-command"
  local shortDescription="A command that outputs some-output"
  writeDummyRunCommand "$commandName"
  writeEmptyJsonObjectAsConfig "$commandName"

  run bats_pipe mrtRun "-h" \| grep "$commandName"

	assert_output "  $commandName Executes run command $commandName"
}

@test "if command config is completely empty, it should exit with an error" {
	local commandName="some-command"
  local shortDescription="A command that outputs some-output"
  local configFile; configFile="$(configFilePath "$commandName")"
  writeDummyRunCommand "$commandName"
  touch "$configFile"

  run mrtRun "-h"

  assert_equal "$status" 1
  assert_line --index 0 "Error while reading $configFile"
  assert_line --index 1 "While parsing config: unexpected end of JSON input"
}