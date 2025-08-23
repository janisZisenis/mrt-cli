mrtSetup() {
  bats_load_library 'mrt/execute.bash'

	mrtExecute setup "$@"
}

mrtSetupAll() {
	mrtSetup all "$@"
}

mrtSetupGitHooks() {
	mrtSetup install-git-hooks
}

mrtSetupClone() {
  mrtSetup clone-repositories
}