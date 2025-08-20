bats_load_library 'ssh-authenticate'
bats_load_library 'common_fixture'
bats_load_library 'repositoriesPath'
bats_load_library 'setup'
bats_load_library 'git'
bats_load_library 'assertLineReversed'

repositoriesPath=$(default_repositories_path)

repositoriesDir() {
	echo "$(testEnvDir)/$repositoriesPath"
}

setup() {
	common_setup
	authenticate
}

teardown() {
	revoke-authentication
	common_teardown
}

@test "if team json does not contain repositoriesPath it clones repository into 'repositories' folder" {
	repositories=("1_TestRepository")

	run setupClone "${repositories[@]}"

	assert_dir_exist "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains an existing repository it should print a messages about successful cloning" {
	repository="1_TestRepository"
	repositoryUrl="$(getRepositoryUrls "$repository")"

	run setupCloneUrls "$repositoryUrl"

	assert_line --index 1 "Cloning $repositoryUrl"
	assert_line --index 3 --regexp "Enumerating objects: [0-9]+, done."
	assert_line_reversed_output 1 "Successfully cloned $repositoryUrl"
}

@test "if team json contains repositoriesPath it clones the repositories into given repositoriesPath folder" {
	repositoriesPath=xyz
	writeRepositoriesPath "$repositoriesPath"
	repository="1_TestRepository"

	run setupClone "$repository"

	assert_dir_exist "$(repositoriesDir)/$repository/.git"
}

@test "if team json contains already existing repositories it clones remaining repositories and skips existing ones" {
	repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	cloneTestingRepositories "$(repositoriesDir)" "${repositories[0]}"

	run setupClone "${repositories[@]}"

	assert_dir_exist "$(repositoriesDir)/${repositories[0]}/.git"
	assert_dir_exist "$(repositoriesDir)/${repositories[1]}/.git"
}

@test "if team json does not contain any repository it does not clone any repository" {
	run setupClone ""

	assert_dir_not_exist "$(repositoriesDir)"
}

@test "if team json contains non-existing repository it should print out a failure message" {
	repository="not-existing"
	repositoryUrl="$(getRepositoryUrls "$repository")"

	run setupCloneUrls "$repositoryUrl"

	assert_output --partial "fatal: Could not read from remote repository."
}

@test "if team json contains non-existing and existing repository it should clone the existing one" {
	repositories=(
		"1_TestRepository"
		"non-existing"
	)

	run setupClone "${repositories[@]}"

	assert_dir_exist "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains repositoriesPrefixes should trim the prefixes while cloning the repositories" {
	repositories=(
		"Prefix1_TestRepository1"
		"Prefix2_TestRepository2"
	)
	writeRepositoriesPrefixes "Prefix1_" "Prefix2_"

	run setupClone "${repositories[@]}"

	assert_dir_exist "$(repositoriesDir)/TestRepository1/.git"
	assert_dir_exist "$(repositoriesDir)/TestRepository2/.git"
}

@test "if team json contains repositoriesPrefixes it should not trim when the prefixes are not in the beginning of the repository names" {
	repositories=(
		"Prefix1_TestRepository1"
		"Prefix2_TestRepository2"
	)
	writeRepositoriesPrefixes "TestRepository1" "TestRepository2"

	run setupClone "${repositories[@]}"

	assert_dir_exist "$(repositoriesDir)/${repositories[0]}/.git"
	assert_dir_exist "$(repositoriesDir)/${repositories[1]}/.git"
}
