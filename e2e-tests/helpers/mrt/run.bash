mrtRun() {
  bats_load_library 'mrt/execute.bash'

	mrtExecute run "$@"
}