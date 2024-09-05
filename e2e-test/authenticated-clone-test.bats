testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'helpers/assertDirectoryExists'
  load 'helpers/assertDirectoryDoesNotExist'
  load 'helpers/writeTeamFile'
  load 'helpers/ssh-authenticate'
  load 'helpers/common'
  load 'helpers/defaults'
  load 'helpers/setupRepositories'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

@test "if team json does not contain repositoriesPath 'setup' clones repository into 'repositories' folder" {
  repository=1_TestRepository

  run setupRepositories "$(testEnvDir)" "$repository"

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/.git"
}

@test "if team json contains repositoriesPath 'setup' clones the repositories into given repositoriesPath folder" {
  repositoriesPath=xyz
  repository=1_TestRepository
  writeRepositoriesPath "$(testEnvDir)" "$repositoriesPath"

  run setupRepositories "$(testEnvDir)" "$repository"

  assert_directory_exists "$(testEnvDir)/$repositoriesPath/$repository/.git"
}

@test "if team json contains already existing repositories 'setup' clones remaining repositories and skips existing ones" {
  firstRepository=1_TestRepository
  secondRepository=2_TestRepository
  git clone git@github-testing:janisZisenisTesting/$firstRepository.git "$(testEnvDir)/$(default_repositories_dir)/$firstRepository"

  run setupRepositories "$(testEnvDir)" "$firstRepository" "$secondRepository"

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/$firstRepository/.git"
  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/$secondRepository/.git"
}

@test "if team json does not contains any repository, 'setup' does not clone any repository" {
  run setupRepositories "$(testEnvDir)" ""

  assert_directory_does_not_exist "$(testEnvDir)/$(default_repositories_dir)"
}

@test "if team json contains non-existing repository, 'setup' should print out a message" {
  nonExistingRepository="not-existing"

  run setupRepositories "$(testEnvDir)" "$nonExistingRepository"

  assert_output "Repository git@github-testing:janisZisenisTesting/$nonExistingRepository.git was not found. Skipping it"
}

@test "if team json contains non-existing and existing repository, 'setup' should clone the existing one" {
  repositoryName=1_TestRepository

  run setupRepositories "$(testEnvDir)" "$repositoryName" "non-exising"

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/$repositoryName/.git"
}

@test "if team json contains repositories but running without 'setup' does not clone the repositories" {
  writeRepositories "$(testEnvDir)" "1_TestRepository"

  run "$(testEnvDir)"/mrt

  assert_directory_does_not_exist "$(testEnvDir)/$(default_repositories_dir)"
}

@test "if team json contains repositoriesPrefixes 'setup' should ignore the prefixes while cloning the repositories" {
  writeRepositoriesPrefixes "$(testEnvDir)" "Prefix1_" "Prefix2_"

  run setupRepositories "$(testEnvDir)" "Prefix1_TestRepository1" "Prefix2_TestRepository2"

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/TestRepository1/.git"
  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/TestRepository2/.git"
}

@test "if team json contains repositoriesPrefixes 'setup' should not ignore the prefixes when the prefixes are not in the beginning of the repository names" {
  writeRepositoriesPrefixes "$(testEnvDir)" "TestRepository1" "TestRepository2"

  run setupRepositories "$(testEnvDir)" "Prefix1_TestRepository1" "Prefix2_TestRepository2"

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/Prefix1_TestRepository1/.git"
  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/Prefix2_TestRepository2/.git"
}
