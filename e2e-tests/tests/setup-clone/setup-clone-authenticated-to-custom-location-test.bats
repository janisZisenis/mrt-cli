bats_load_library 'fixtures/authenticated_fixture.bash'
bats_load_library 'fixtures/common_fixture.bash'
bats_load_library 'mrt/clone.bash'
bats_load_library 'writeTeamFile.bash'

setup() {
	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "if team json contains repositoriesPath it clones the repositories into given repositoriesPath folder" {
	repositoriesPath="xyz"
	writeRepositoriesPath "$repositoriesPath"
	repository="1_TestRepository"

	run clone_repositories_using_mrt "$repository"

	assert_dir_exist "$(testEnvDir)/$repositoriesPath/$repository/.git"
}
