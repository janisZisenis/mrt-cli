bats_load_library 'setup'
bats_load_library 'repositoriesPath'
bats_load_library 'testRepositories'
bats_load_library "fixtures/authenticated_fixture"

setup() {
	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "if setup is run with skipping the clone step it should not clone the repositories" {
	repository="1_TestRepository"
	repositoryUrl="$(getTestingRepositoryUrl "$repository")"
	writeRepositoriesUrls "$repositoryUrl"

	run execute setup all --skip-clone-repositories

	assert_dir_not_exist "$(testEnvDir)/$(default_repositories_path)/$repository"
}

@test "if setup is run with skipping the clone step it should print a skip message" {
	repository="1_TestRepository"
	repositoryUrl="$(getTestingRepositoryUrl "$repository")"
	writeRepositoriesUrls "$repositoryUrl"

	run execute setup all --skip-clone-repositories

	assert_line --index 0 "Skipping clone-repositories step."
}
