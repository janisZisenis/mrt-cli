write_script() {
  local content="$1"
  local script_path="$2"

	mkdir -p "$(dirname "$script_path")"

  echo "$content" > "$script_path"
  chmod +x "$script_path"
}

write_dummy_script() {
	local script_path="$1"
	local exit_code="$2"
	local output="$3"

  bats_load_library 'scripts/script_factory.bash'
  write_script "$(make_dummy_script)" "$script_path"
}

write_stub_script() {
	local script_path="$1"
	local exit_code="$2"
	local output="$3"

  bats_load_library 'scripts/script_factory.bash'
  write_script "$(make_stub_script "$output" "$exit_code")" "$script_path"
}

write_spy_script() {
	local script_path="$1"

  local command_name; command_name="$(basename "$script_path")"

  bats_load_library 'scripts/script_factory.bash'
  local content; content="$(make_spy_script "$command_name")"

	write_script "$content" "$script_path"
}

assert_script_was_executed_with_parameters() {
	local script_path="$1"
	local expected_parameters="$2"

  assert_script_was_executed "$1"

  bats_load_library 'scripts/script_factory.bash'
	assert_equal "$(cat "$script_path$(spy_file_suffix)")" "$expected_parameters"
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
	local script_path=$1

  bats_load_library 'scripts/script_factory.bash'
  write_script "$(make_script_requesting_input)" "$script_path"
}

assert_script_received_input() {
  local script_path="$1"
  local expected_input="$2"

	assert_file_exist "$(dirname "$script_path")/$expected_input"
}

write_std_err_script() {
	local script_path=$1
	local error_message=$2

  bats_load_library 'scripts/script_factory.bash'
	write_script "$(make_std_err_script "$error_message")" "$script_path"
}