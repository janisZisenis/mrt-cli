load '../../helpers/directoryAssertions'
load '../../helpers/directoryAssertions'
load '../../helpers/ssh-authenticate'
load '../../helpers/common'
load '../../helpers/repositoriesPath'
load '../../helpers/setup'
load '../../helpers/git'
load '../../helpers/assertLineReversed'

repositoriesPath=$(default_repositories_path)

repositoriesDir() {
	echo "$_testEnvDir/$repositoriesPath"
}

setup() {
	_common_setup
	authenticate
}

teardown() {
	revoke-authentication
	_common_teardown
}

@test "if team json does not contain repositoriesPath it clones repository into 'repositories' folder" {
	repositories=("1_TestRepository")

	run setupClone "${repositories[@]}"

	assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains an existing repository it should print a messages about successful cloning" {
	repository="1_TestRepository"
	repositoryUrl="$(getRepositoryUrls "$repository")"

	run setupCloneUrls "$repositoryUrl"

	assert_line --index 1 "Cloning $repositoryUrl into $repositoriesPath/$repository"
	assert_line --index 2 --regexp "Enumerating objects: [0-9]+, done."
	assert_line_reversed_output 1 "Successfully cloned $repositoryUrl"
}

@test "if team json contains repositoriesPath it clones the repositories into given repositoriesPath folder" {
	repositoriesPath=xyz
	writeRepositoriesPath "$repositoriesPath"
	repository="1_TestRepository"

	run setupClone "$repository"

	assert_directory_exists "$(repositoriesDir)/$repository/.git"
}

@test "if team json contains already existing repositories it clones remaining repositories and skips existing ones" {
	repositories=(
		"1_TestRepository"
		"2_TestRepository"
	)
	cloneTestingRepositories "$(repositoriesDir)" "${repositories[0]}"

	run setupClone "${repositories[@]}"

	assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
	assert_directory_exists "$(repositoriesDir)/${repositories[1]}/.git"
}

@test "if team json does not contain any repository it does not clone any repository" {
	run setupClone ""

	assert_directory_does_not_exist "$(repositoriesDir)"
}

@test "if team json contains non-existing repository it should print out a failure message" {
	repository="not-existing"
	repositoryUrl="$(getRepositoryUrls "$repository")"

	run setupCloneUrls "$repositoryUrl"

	assert_line --index 1 "Cloning $repositoryUrl into $(default_repositories_path)/$repository"
	assert_line --index 2 "Repository $repositoryUrl was not found. Skipping it"
}

@test "if team json contains non-existing and existing repository it should clone the existing one" {
	repositories=(
		"1_TestRepository"
		"non-existing"
	)

	run setupClone "${repositories[@]}"

	assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains repositoriesPrefixes should trim the prefixes while cloning the repositories" {
	repositories=(
		"Prefix1_TestRepository1"
		"Prefix2_TestRepository2"
	)
	writeRepositoriesPrefixes "Prefix1_" "Prefix2_"

	run setupClone "${repositories[@]}"

	assert_directory_exists "$(repositoriesDir)/TestRepository1/.git"
	assert_directory_exists "$(repositoriesDir)/TestRepository2/.git"
}

@test "if team json contains repositoriesPrefixes it should not trim when the prefixes are not in the beginning of the repository names" {
	repositories=(
		"Prefix1_TestRepository1"
		"Prefix2_TestRepository2"
	)
	writeRepositoriesPrefixes "TestRepository1" "TestRepository2"

	run setupClone "${repositories[@]}"

	assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
	assert_directory_exists "$(repositoriesDir)/${repositories[1]}/.git"
}
