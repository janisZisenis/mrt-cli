bats_load_library 'mrt/setup'
bats_load_library 'repositoriesPath'
bats_load_library 'writeTeamFile'
bats_load_library 'git'
bats_load_library 'testRepositories'
bats_load_library "fixtures/authenticated_fixture"

setup() {
  authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "If setup is run with skipping git hooks, it should print skip message" {
  run mrtSetupAll --skip-install-git-hooks

	assert_line --partial "Skipping install-git-hooks step."
}
