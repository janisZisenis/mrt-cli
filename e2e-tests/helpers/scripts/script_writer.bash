write_script() {
  local content="$1"
  local scriptPath="$2"

	mkdir -p "$(dirname "$scriptPath")"

  echo "$content" > "$scriptPath"
  chmod +x "$scriptPath"
}

write_dummy_script() {
	local scriptPath="$1"
	local exitCode="$2"
	local output="$3"

  bats_load_library 'scripts/script_factory.bash'
  write_script "$(make_dummy_script)" "$scriptPath"
}

write_stub_script() {
	local scriptPath="$1"
	local exitCode="$2"
	local output="$3"

  bats_load_library 'scripts/script_factory.bash'
  write_script "$(make_stub_script "$output" "$exitCode")" "$scriptPath"
}

write_spy_script() {
	local scriptPath="$1"

  local commandName; commandName="$(basename "$scriptPath")"

  bats_load_library 'scripts/script_factory.bash'
  local content; content="$(make_spy_script "$commandName")"

	write_script "$content" "$scriptPath"
}

assert_script_was_executed_with_parameters() {
	local scriptPath="$1"
	local expectedParameters="$2"

  assert_script_was_executed "$1"

  bats_load_library 'scripts/script_factory.bash'
	assert_equal "$(cat "$scriptPath$(spy_file_suffix)")" "$expectedParameters"
}

assert_script_was_executed() {
  bats_load_library 'scripts/script_factory.bash'
	assert_file_exist "$1$(spy_file_suffix)"
}

assert_script_was_not_executed() {
  bats_load_library 'scripts/script_factory.bash'
	assert_file_not_exist "$1$(spy_file_suffix)"
}

write_script_requesting_input() {
	local scriptPath=$1

  bats_load_library 'scripts/script_factory.bash'
  write_script "$(make_script_requesting_input)" "$scriptPath"
}

assert_script_received_input() {
  local scriptPath="$1"
  local expectedInput="$2"

	assert_file_exist "$(dirname "$scriptPath")/$expectedInput"
}

write_std_err_script() {
	local scriptPath=$1
	local errorMessage=$2

  bats_load_library 'scripts/script_factory.bash'
	write_script "$(make_std_err_script "$errorMessage")" "$scriptPath"
}