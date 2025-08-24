setup() {
  bats_load_library 'fixtures/common_fixture.bash'
  bats_load_library 'mrt/setup.bash'
  bats_load_library 'mrt/clone.bash'
  bats_load_library 'write_team_file.bash'
  bats_load_library 'assert_line_reversed.bash'
  bats_load_library 'test_repositories.bash'

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
	local repository_path="$1"
	local repositories=("1_TestRepository")
	write_repositories_path "$repository_path"

	run clone_repositories_using_mrt "${repositories[@]}"

	assert_line --index 0 "Start cloning repositories into \"$repository_path\""
	assert_line --index 9 "Cloning repositories done"
}

@test "If team json contains 2 repositories it should print out a done message after cloning second" {
	local repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)

	run clone_repositories_using_mrt "${repositories[@]}"

	assert_line_reversed_output 0 "Cloning repositories done"
}

@test "if team json contains existing repositories but authentication is missing it should print a failure message" {
	local repository="1_TestRepository"
	local repository_url; repository_url="$(get_testing_repository_url "$repository")"

	run clone_repository_urls_using_mrt "$repository_url"

	assert_line --index 1 "Cloning $repository_url"
	assert_line --index 7 "Failed to clone repository, skipping it."
}

@test "if team json does not contain any repositories it prints out a message" {
	local repositories_urls=()
	write_repositories_urls "${repositories_urls[@]}"

	run mrt_setup_clone

	assert_success
	assert_output 'The team file does not contain any repositories, no repositories to clone.'
}

@test "if team json does not exist it prints out a message" {
	run mrt_setup_clone

	assert_success
	assert_output 'Could not read team file. To setup your repositories create a "team.json" file and add repositories to it.'
}
