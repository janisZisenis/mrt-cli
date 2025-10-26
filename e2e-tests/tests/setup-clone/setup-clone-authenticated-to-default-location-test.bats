repositories_dir() {
	bats_load_library 'fixtures/common_fixture.bash'
	bats_load_library 'repositories_path.bash'

	echo "$(test_env_dir)/$(default_repositories_path)"
}

cloned_git_folder() {
	local repository="$1"

	echo "$(repositories_dir)/$repository/.git"
}

setup() {
	bats_load_library 'fixtures/authenticated_fixture.bash'
	bats_load_library 'test_repositories.bash'
	bats_load_library 'mrt/clone.bash'
	bats_load_library 'git.bash'
	bats_load_library 'assert_line_reversed.bash'
	bats_load_library 'write_team_file.bash'

	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "if team json does not contain repositoriesPath it clones repository into 'repositories' folder" {
	local repositories=("1_TestRepository")

	run clone_repositories_using_mrt "${repositories[@]}"

	assert_dir_exist "$(cloned_git_folder "${repositories[0]}")"
}

@test "if team json contains an existing repository it should print a messages about successful cloning" {
	local repository="1_TestRepository"
	local repository_url
	repository_url="$(get_repository_urls "$repository")"

	run clone_repository_urls_using_mrt "$repository_url"

	assert_line --index 1 "Cloning $repository_url"
	assert_line --index 3 --regexp "Enumerating objects: [0-9]+, done."
	assert_line_reversed_output 1 "Successfully cloned $repository_url"
}

@test "if team json contains already existing repositories it clones remaining repositories and skips existing ones" {
	local repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	clone_testing_repositories "$(repositories_dir)" "${repositories[0]}"

	run clone_repositories_using_mrt "${repositories[@]}"

	assert_dir_exist "$(cloned_git_folder "${repositories[0]}")"
	assert_dir_exist "$(cloned_git_folder "${repositories[1]}")"
}

@test "if team json does not contain any repository it does not clone any repository" {
	local no_repositories=()

	run clone_repositories_using_mrt "${no_repositories[@]}"

	assert_dir_not_exist "$(repositories_dir)"
}

@test "if team json contains non-existing repository it should print out a failure message" {
	local repository="not-existing"
	local repository_url
	repository_url="$(get_repository_urls "$repository")"

	run clone_repository_urls_using_mrt "$repository_url"

	assert_output --partial "fatal: Could not read from remote repository."
}

@test "if team json contains non-existing and existing repository it should clone the existing one" {
	local repositories=(
		"1_TestRepository"
		"non-existing"
	)

	run clone_repositories_using_mrt "${repositories[@]}"

	assert_dir_exist "$(cloned_git_folder "${repositories[0]}")"
}

@test "if team json contains repositoriesPrefixes should trim the prefixes while cloning the repositories" {
	local repositories=(
		"Prefix1_TestRepository1"
		"Prefix2_TestRepository2"
	)
	write_repositories_prefixes "Prefix1_" "Prefix2_"

	run clone_repositories_using_mrt "${repositories[@]}"

	assert_dir_exist "$(cloned_git_folder "TestRepository1")"
	assert_dir_exist "$(cloned_git_folder "TestRepository2")"
}

@test "if team json contains repositoriesPrefixes it should not trim when the prefixes are not in the beginning of the repository names" {
	local repositories=(
		"Prefix1_TestRepository1"
		"Prefix2_TestRepository2"
	)
	write_repositories_prefixes "TestRepository1" "TestRepository2"

	run clone_repositories_using_mrt "${repositories[@]}"

	assert_dir_exist "$(cloned_git_folder "${repositories[0]}")"
	assert_dir_exist "$(cloned_git_folder "${repositories[1]}")"
}
