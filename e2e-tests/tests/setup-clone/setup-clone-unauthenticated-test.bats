bats_load_library 'common_fixture'
bats_load_library 'common_fixture'
bats_load_library 'repositoriesPath'
bats_load_library 'setup'
bats_load_library 'writeTeamFile'
bats_load_library 'assertLineReversed'
bats_load_library 'testRepositories'

setup() {
	common_setup
}

teardown() {
	common_teardown
}

@test "If team json contains repository and some repository path it should print out message, that it clones the repositories" {
	test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories "some-path"
}

@test "If team json contains repository and another repository path it should print out message, that it clones the repositories" {
	test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories "another-path"
}

test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories() {
	repositoryPath=$1
	repositories=(
		"1_TestRepository"
	)
	writeRepositoriesPath "$repositoryPath"

	run setupClone "${repositories[@]}"

	assert_line --index 0 "Start cloning repositories into \"$repositoryPath\""
	assert_line --index 9 "Cloning repositories done"
}

@test "If team json contains 2 repositories it should print out a done message after cloning second" {
	repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)

	run setupClone "${repositories[@]}"

	assert_line_reversed_output 0 "Cloning repositories done"
}

@test "if team json contains existing repositories but authentication is missing it should print a failure message" {
	repository="1_TestRepository"
	repositoryUrl="$(getTestingRepositoryUrl "$repository")"

	run setupCloneUrls "$repositoryUrl"

	assert_line --index 1 "Cloning $repositoryUrl"
	assert_line --index 7 "Failed to clone repository, skipping it."
}

@test "if team json does not contain any repositories it prints out a message" {
	repositoriesUrls=()
	writeRepositoriesUrls "${repositoriesUrls[@]}"

	run execute setup clone-repositories

	assert_success
	assert_output 'The team file does not contain any repositories, no repositories to clone.'
}

@test "if team json does not exist it prints out a message" {
	run execute setup clone-repositories

	assert_success
	assert_output 'Could not read team file. To setup your repositories create a "team.json" file and add repositories to it.'
}
