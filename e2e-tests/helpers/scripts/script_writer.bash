writeScript() {
  local content="$1"
  local scriptPath="$2"

	mkdir -p "$(dirname "$scriptPath")"

  echo "$content" > "$scriptPath"
  chmod +x "$scriptPath"
}

writeDummyScript() {
	local scriptPath="$1"
	local exitCode="$2"
	local output="$3"

  bats_load_library 'scripts/script_factory.bash'
  writeScript "$(makeDummyScript)" "$scriptPath"
}

writeStubScript() {
	local scriptPath="$1"
	local exitCode="$2"
	local output="$3"

  bats_load_library 'scripts/script_factory.bash'
  writeScript "$(makeStubScript "$output" "$exitCode")" "$scriptPath"
}

writeSpyScript() {
	local scriptPath="$1"

  local commandName
  commandName="$(basename "$scriptPath")"

  bats_load_library 'scripts/script_factory.bash'
  local content
  content="$(makeSpyScript "$commandName")"

	writeScript "$content" "$scriptPath"
}

assert_script_was_executed_with_parameters() {
	local scriptPath="$1"
	local expectedParameters="$2"

  assert_script_was_executed "$1"

  bats_load_library 'scripts/script_factory.bash'
	assert_equal "$(cat "$scriptPath$(spyFileSuffix)")" "$expectedParameters"
}

assert_script_was_executed() {
  bats_load_library 'scripts/script_factory.bash'
	assert_file_exist "$1$(spyFileSuffix)"
}

assert_script_was_not_executed() {
  bats_load_library 'scripts/script_factory.bash'
	assert_file_not_exist "$1$(spyFileSuffix)"
}

writeScriptRequestingInput() {
	local scriptPath=$1

  bats_load_library 'scripts/script_factory.bash'
  writeScript "$(makeScriptRequestingInput)" "$scriptPath"
}

assert_script_received_input() {
  local scriptPath="$1"
  local expectedInput="$2"

	assert_file_exist "$(dirname "$scriptPath")/$expectedInput"
}

writeStdErrScript() {
	local scriptPath=$1
	local errorMessage=$2

  bats_load_library 'scripts/script_factory.bash'
	writeScript "$(makeStdErrScript "$errorMessage")" "$scriptPath"
}