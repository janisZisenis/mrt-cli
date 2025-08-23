assert_line_reversed_output() {
	local index="$1"
	local expected="$2"

	if ! [[ "$index" =~ ^-?[0-9]+$ ]]; then
		echo "Error: the given index is not a valid integer."
		return 1
	fi

  # we want a zero-based notation
  # but using 'tail' it would be is one-based
  local zero_based_index=$((index + 1))

	local actual
	# shellcheck disable=SC2154
	# output is populated by bats
	actual=$(echo "$output" | tail -n $zero_based_index | head -n 1)
	assert_equal "$actual" "$expected"
}
