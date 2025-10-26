setup() {
	bats_load_library 'mrt/setup.bash'
	bats_load_library 'fixtures/authenticated_fixture.bash'

	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "If setup is run with skipping git hooks, it should print skip message" {
	run mrt_setup_all --skip-install-git-hooks

	assert_line --partial "Skipping install-git-hooks step."
}
