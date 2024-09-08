load 'helpers/assertDirectoryExists'
load 'helpers/assertDirectoryDoesNotExist'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/setupRepositories'

testEnvDir=$(_testEnvDir)

setup() {
  _common_setup "$testEnvDir"
  authenticate
}

teardown() {
  _common_teardown "$testEnvDir"
  revoke-authentication
}

@test "if team json does not contain repositoriesPath 'setup all' clones repository into 'repositories' folder" {
  repository=1_TestRepository
  repositoryUrl=$(getTestingRepositoryUrl $repository)

  run setupRepositories "$testEnvDir" "$repositoryUrl"

  assert_directory_exists "$testEnvDir/$(default_repositories_dir)/$repository/.git"
}

@test "if team json contains an existing repository 'setup all' should print a messages about successful cloning" {
  repository=1_TestRepository
  repositoryUrl=$(getTestingRepositoryUrl $repository)

  run setupRepositories "$testEnvDir" "$repositoryUrl"

  assert_line --index 1 "Cloning git@github-testing:janisZisenisTesting/$repository.git into $(default_repositories_dir)/$repository"
  assert_line --index 2 "Successfully cloned git@github-testing:janisZisenisTesting/$repository.git"
}

@test "if team json contains repositoriesPath 'setup all' clones the repositories into given repositoriesPath folder" {
  repositoriesPath=xyz
  repository=1_TestRepository
  repositoryUrl=$(getTestingRepositoryUrl $repository)
  writeRepositoriesPath "$testEnvDir" "$repositoriesPath"

  run setupRepositories "$testEnvDir" "$repositoryUrl"

  assert_directory_exists "$testEnvDir/$repositoriesPath/$repository/.git"
}

@test "if team json contains already existing repositories 'setup all' clones remaining repositories and skips existing ones" {
  firstRepository=1_TestRepository
  firstRepositoryUrl=$(getTestingRepositoryUrl "$firstRepository")
  secondRepository=2_TestRepository
  secondRepositoryUrl=$(getTestingRepositoryUrl "$secondRepository")
  git clone "$firstRepositoryUrl" "$testEnvDir/$(default_repositories_dir)/$firstRepository"

  run setupRepositories "$testEnvDir" "$firstRepositoryUrl" "$secondRepositoryUrl"

  assert_directory_exists "$testEnvDir/$(default_repositories_dir)/$firstRepository/.git"
  assert_directory_exists "$testEnvDir/$(default_repositories_dir)/$secondRepository/.git"
}

@test "if team json does not contains any repository, 'setup all' does not clone any repository" {
  run setupRepositories "$testEnvDir" ""

  assert_directory_does_not_exist "$testEnvDir/$(default_repositories_dir)"
}

@test "if team json contains non-existing repository, 'setup all' should print out a failure message" {
  nonExistingRepository="not-existing"
  repositoryUrl=$(getTestingRepositoryUrl "$nonExistingRepository")

  run setupRepositories "$testEnvDir" "$repositoryUrl"

  assert_line --index 1 "Cloning $repositoryUrl into $(default_repositories_dir)/$nonExistingRepository"
  assert_line --index 2 "Repository $repositoryUrl was not found. Skipping it"
}

@test "if team json contains non-existing and existing repository, 'setup all' should clone the existing one" {
  repositoryName=1_TestRepository
  repositoryUrl=$(getTestingRepositoryUrl "$repositoryName")

  run setupRepositories "$testEnvDir" "$repositoryUrl" "$(getTestingRepositoryUrl "non-existing")"

  assert_directory_exists "$testEnvDir/$(default_repositories_dir)/$repositoryName/.git"
}

@test "if team json contains repositories but running without 'setup all' does not clone the repositories" {
  writeRepositories "$testEnvDir" "$(getTestingRepositoryUrl "1_TestRepository")"

  run "$testEnvDir"/mrt

  assert_directory_does_not_exist "$testEnvDir/$(default_repositories_dir)"
}

@test "if team json contains repositoriesPrefixes 'setup all' should ignore the prefixes while cloning the repositories" {
  writeRepositoriesPrefixes "$testEnvDir" "Prefix1_" "Prefix2_"

  run setupRepositories "$testEnvDir" "$(getTestingRepositoryUrl "Prefix1_TestRepository1")" "$(getTestingRepositoryUrl "Prefix2_TestRepository2")"

  assert_directory_exists "$testEnvDir/$(default_repositories_dir)/TestRepository1/.git"
  assert_directory_exists "$testEnvDir/$(default_repositories_dir)/TestRepository2/.git"
}

@test "if team json contains repositoriesPrefixes 'setup all' should not ignore the prefixes when the prefixes are not in the beginning of the repository names" {
  writeRepositoriesPrefixes "$testEnvDir" "TestRepository1" "TestRepository2"

  run setupRepositories "$testEnvDir" "$(getTestingRepositoryUrl "Prefix1_TestRepository1")" "$(getTestingRepositoryUrl "Prefix2_TestRepository2")"

  assert_directory_exists "$testEnvDir/$(default_repositories_dir)/Prefix1_TestRepository1/.git"
  assert_directory_exists "$testEnvDir/$(default_repositories_dir)/Prefix2_TestRepository2/.git"
}
