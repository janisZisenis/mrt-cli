test_env_dir() {
	echo "$BATS_TEST_TMPDIR"
}

common_setup() {
	bats_load_library bats-assert
	bats_load_library bats-support
	bats_load_library bats-file

	PATH=$PATH:"$(mrt run binary-location -- --dir)"
	eval "$(ssh-agent -s 3>&-)"
	export MRT_BASE_SSH_COMMAND="$GIT_SSH_COMMAND"
	export GIT_SSH_COMMAND="$GIT_SSH_COMMAND -o PubkeyAuthentication=no"
}

common_teardown() {
	eval "$(ssh-agent -k)"
}
