bats_load_library 'mrt/setup.bash'
bats_load_library 'repositoriesPath.bash'
bats_load_library 'writeTeamFile.bash'
bats_load_library 'git.bash'
bats_load_library 'testRepositories.bash'
bats_load_library 'fixtures/authenticated_fixture.bash'

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
