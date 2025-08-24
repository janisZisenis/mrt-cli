mrt_setup() {
  bats_load_library 'mrt/execute.bash'

	mrt_execute setup "$@"
}

mrt_setup_all() {
	mrt_setup all "$@"
}

mrt_setup_git_hooks() {
	mrt_setup install-git-hooks
}

mrt_setup_clone() {
  mrt_setup clone-repositories
}