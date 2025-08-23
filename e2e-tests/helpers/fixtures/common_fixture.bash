testEnvDir() {
	echo "$BATS_TEST_TMPDIR"
}

common_setup() {
	bats_load_library bats-assert
	bats_load_library bats-support
  bats_load_library bats-file

	PATH=$PATH:"$(mrt run binary-location -- --dir)"
	eval "$(ssh-agent -s 3>&-)"
}

common_teardown() {
	eval "$(ssh-agent -k)"
}