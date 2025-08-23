writeDummyRunCommand() {
	local commandName="$1"

  bats_load_library 'commands/commandWriter.bash'
  bats_load_library 'commands/run/runCommandLocation.bash'
  writeDummyCommand "$(runCommandLocation)" "$commandName"
}

writeStubRunCommand() {
	local commandName="$1"
	local exitCode="$2"
	local output="$3"

  bats_load_library 'commands/commandWriter.bash'
  bats_load_library 'commands/run/runCommandLocation.bash'
  writeStubCommand "$(runCommandLocation)" "$commandName" "$exitCode" "$output"
}

writeSpyRunCommand() {
	local commandName="$1"

  bats_load_library 'commands/commandWriter.bash'
  bats_load_library 'commands/run/runCommandLocation.bash'
  writeSpyCommand "$(runCommandLocation)" "$commandName"
}

assert_run_command_was_executed_with_parameters() {
	local commandName="$1"
	local expectedParameters="$2"

  bats_load_library 'commands/commandWriter.bash'
  bats_load_library 'commands/run/runCommandLocation.bash'
  assert_command_was_executed_with_parameters "$(runCommandLocation)" "$commandName" "$expectedParameters"
}

writeStdErrRunCommand() {
	local commandName="$1"
	local error="$2"

  bats_load_library 'commands/commandWriter.bash'
	bats_load_library 'commands/run/runCommandLocation.bash'
	writeStdErrCommand "$(runCommandLocation)" "$commandName" "$error"
}

writeRunCommandRequestingInput() {
	local commandName="$1"

  bats_load_library 'commands/commandWriter.bash'
	bats_load_library 'commands/run/runCommandLocation.bash'
	writeCommandRequestingInput "$(runCommandLocation)" "$commandName"
}

assert_run_command_received_input() {
	local commandName="$1"
	local expectedInput="$2"

  bats_load_library 'commands/commandWriter.bash'
	bats_load_library 'commands/run/runCommandLocation.bash'
	assert_command_received_input "$(runCommandLocation)" "$commandName" "$expectedInput"
}