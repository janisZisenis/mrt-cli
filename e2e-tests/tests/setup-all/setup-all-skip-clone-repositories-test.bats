setup() {
  bats_load_library 'repositoriesPath.bash'
  bats_load_library 'testRepositories.bash'
  bats_load_library 'fixtures/authenticated_fixture.bash'
  bats_load_library 'mrt/setup.bash'
  bats_load_library 'writeTeamFile.bash'

	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "if setup is run with skipping the clone step it should not clone the repositories" {
	repository="1_TestRepository"
	repositoryUrl="$(getTestingRepositoryUrl "$repository")"
	writeRepositoriesUrls "$repositoryUrl"

	run mrtSetupAll --skip-clone-repositories

	assert_dir_not_exist "$(testEnvDir)/$(default_repositories_path)/$repository"
}

@test "if setup is run with skipping the clone step it should print a skip message" {
	repository="1_TestRepository"
	repositoryUrl="$(getTestingRepositoryUrl "$repository")"
	writeRepositoriesUrls "$repositoryUrl"

	run mrtSetupAll --skip-clone-repositories

	assert_line --index 0 "Skipping clone-repositories step."
}
