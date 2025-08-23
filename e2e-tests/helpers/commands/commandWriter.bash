commandFileName() {
	echo "command"
}

writeDummyCommand() {
  local commandLocation="$1"
  local commandName="$2"

  bats_load_library "scripts/scriptWriter.bash"
  writeDummyScript "$commandLocation/$commandName/$(commandFileName)"
}

writeStubCommand() {
	local commandLocation="$1"
	local commandName="$2"
	local exitCode="$3"
	local output="$4"

  bats_load_library "scripts/scriptWriter.bash"
	writeStubScript "$commandLocation/$commandName/$(commandFileName)" "$exitCode" "$output"
}

writeSpyCommand() {
	local commandLocation="$1"
	local commandName="$2"

  bats_load_library "scripts/scriptWriter.bash"
	writeSpyScript "$commandLocation/$commandName/$(commandFileName)"
}

assert_command_was_executed_with_parameters() {
	local commandLocation="$1"
	local commandName="$2"
	local expectedParameters="$3"

  bats_load_library "scripts/scriptWriter.bash"
	assert_script_was_executed_with_parameters "$commandLocation/$commandName/$(commandFileName)" "$expectedParameters"
}

assert_command_was_not_executed() {
	local commandLocation="$1"
	local commandName="$2"

  bats_load_library "scripts/scriptWriter.bash"
	assert_script_was_not_executed "$commandLocation/$commandName/$(commandFileName)"
}

writeStdErrCommand() {
	local commandLocation="$1"
	local commandName="$2"
	local error="$3"

  bats_load_library "scripts/scriptWriter.bash"
	writeStdErrScript "$commandLocation/$commandName/$(commandFileName)" "$error"
}

writeCommandRequestingInput() {
	local commandLocation="$1"
	local commandName="$2"

  bats_load_library "scripts/scriptWriter.bash"
	writeScriptRequestingInput "$commandLocation/$commandName/$(commandFileName)"
}

assert_command_received_input() {
	local commandLocation="$1"
	local commandName="$2"
	local expectedInput="$3"

  bats_load_library "scripts/scriptWriter.bash"
  assert_script_received_input "$commandLocation/$commandName/$(commandFileName)" "$expectedInput"
}
