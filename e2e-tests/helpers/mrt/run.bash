mrt_run() {
	bats_load_library 'mrt/execute.bash'

	mrt_execute run "$@"
}
