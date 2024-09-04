testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'test_helper/assertDirectoryExists'
  load 'test_helper/assertDirectoryDoesNotExist'
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'
  load 'test_helper/defaults'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

@test "if team json does not contain repositoriesPath 'setup' clones repository into 'repositories' folder" {
  repository=1_TestRepository
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/.git"
}

@test "if team json contains repositoriesPath 'setup' clones the repositories into given repositoriesPath folder" {
  repositoriesPath=xyz
  repository=1_TestRepository
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$repositoriesPath/$repository/.git"
}

@test "if team json contains already existing repositories 'setup' clones remaining repositories and skips existing ones" {
  firstRepository=1_TestRepository
  secondRepository=2_TestRepository
  git clone git@github-testing:janisZisenisTesting/$firstRepository.git "$(testEnvDir)/$(default_repositories_dir)/$firstRepository"
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$firstRepository.git\",
          \"git@github-testing:janisZisenisTesting/$secondRepository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/$firstRepository/.git"
  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/$secondRepository/.git"
}

@test "if team json does not contains any repository, 'setup' does not clone any repository" {
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": []
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_does_not_exist "$(testEnvDir)/$(default_repositories_dir)"
}

@test "if team json contains non-existing repository, 'setup' should print out a message" {
  nonExistingRepository=git@github-testing:janisZisenisTesting/not-existing.git
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": [
        \"$nonExistingRepository\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_output "Repository $nonExistingRepository was not found. Skipping it"
}

@test "if team json contains non-existing and existing repository, 'setup' should clone the existing one" {
  repositoryName=1_TestRepository
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/non-exising.git\",
          \"git@github-testing:janisZisenisTesting/$repositoryName.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/$repositoryName/.git"
}

@test "if team json contains repositories but running without 'setup' does not clone the repositories" {
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/1_TestRepository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt

  assert_directory_does_not_exist "$(testEnvDir)/$(default_repositories_dir)"
}

@test "if team json contains repositoriesPrefixes 'setup' should ignore the prefixes while cloning the repositories" {
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPrefixes\": [
        \"Prefix1_\",
        \"Prefix2_\"
      ],
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/Prefix1_TestRepository1.git\",
          \"git@github-testing:janisZisenisTesting/Prefix2_TestRepository2.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/TestRepository1/.git"
  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/TestRepository2/.git"
}

@test "if team json contains repositoriesPrefixes 'setup' should not ignore the prefixes when the prefixes are not in the beginning of the repository names" {
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPrefixes\": [
        \"TestRepository1\",
        \"TestRepository2\"
      ],
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/Prefix1_TestRepository1.git\",
          \"git@github-testing:janisZisenisTesting/Prefix2_TestRepository2.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/Prefix1_TestRepository1/.git"
  assert_directory_exists "$(testEnvDir)/$(default_repositories_dir)/Prefix2_TestRepository2/.git"
}

