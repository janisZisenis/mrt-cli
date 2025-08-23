_setupCommandLocation() {
  bats_load_library 'fixtures/common_fixture.bash'

  echo "$(testEnvDir)/setup"
}

writeStubSetupCommand() {
	local commandName="$1"
	local exitCode="$2"
	local output="$3"

  bats_load_library 'commands/commandWriter.bash'
  writeStubCommand "$(_setupCommandLocation)" "$commandName" "$exitCode" "$output"
}

writeSpySetupCommand() {
	local commandName="$1"

  bats_load_library 'commands/commandWriter.bash'
  writeSpyCommand "$(_setupCommandLocation)" "$commandName"
}

assert_setup_command_was_executed() {
	local commandName="$1"
	local expectedParameters="$2"

  bats_load_library 'commands/commandWriter.bash'
	assert_command_was_executed_with_parameters "$(_setupCommandLocation)" "$commandName" "$expectedParameters"
}

assert_setup_command_was_not_executed() {
	local commandName="$1"

  bats_load_library 'commands/commandWriter.bash'
	assert_command_was_not_executed "$(_setupCommandLocation)" "$commandName"
}

writeStdErrSetupCommand() {
	local commandName="$1"
	local error="$2"

  bats_load_library 'commands/commandWriter.bash'
	writeStdErrCommand "$(_setupCommandLocation)" "$commandName" "$error"
}

writeSetupCommandRequestingInput() {
	local commandName="$1"

  bats_load_library 'commands/commandWriter.bash'
	writeCommandRequestingInput "$(_setupCommandLocation)" "$commandName"
}

assert_setup_command_received_input() {
	local commandName="$1"
	local expectedInput="$2"

  bats_load_library 'commands/commandWriter.bash'
	assert_command_received_input "$(_setupCommandLocation)" "$commandName" "$expectedInput"
}