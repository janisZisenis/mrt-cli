setup() {
	bats_load_library 'fixtures/authenticated_fixture.bash'
	bats_load_library 'fixtures/common_fixture.bash'
	bats_load_library 'mrt/clone.bash'
	bats_load_library 'write_team_file.bash'

	authenticated_setup
}

teardown() {
	authenticated_teardown
}
