testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'test_helper/assertDirectoryExists'
  load 'test_helper/assertDirectoryDoesNotExist'
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

@test "if team json contains repositories 'setup' clones that repository into given repository path" {
  repositoriesPath=repositories
  firstRepository=1_TestRepository
  secondRepository=2_TestRepository
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$firstRepository.git\",
          \"git@github-testing:janisZisenisTesting/$secondRepository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$repositoriesPath/$firstRepository/.git"
  assert_directory_exists "$(testEnvDir)/$repositoriesPath/$secondRepository/.git"
}

@test "if team json contains xyz as repositoriesPath 'setup' clones the repositories into given xyz folder" {
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

@test "if team json contains already existing repositories 'setup' clones remaining repositories given repository path" {
  repositoriesPath=repositories
  firstRepository=1_TestRepository
  secondRepository=2_TestRepository
  git clone git@github-testing:janisZisenisTesting/$firstRepository.git "$(testEnvDir)"/$repositoriesPath/$firstRepository
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$firstRepository.git\",
          \"git@github-testing:janisZisenisTesting/$secondRepository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$repositoriesPath/$firstRepository/.git"
  assert_directory_exists "$(testEnvDir)/$repositoriesPath/$secondRepository/.git"
}

@test "if team json does not contains any repository, 'setup' does not clone any repository" {
  repositoriesPath=repositories
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": []
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_does_not_exist "$(testEnvDir)/$repositoriesPath"
}

@test "if team json does not contains any repository, 'setup' exits with error" {
  repositoriesPath=repositories
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": []
  }"

  run "$(testEnvDir)"/mrt setup

  assert_failure
}

@test "if team json contains non-existing repository, 'setup' should print out a message" {
  repositoriesPath=repositories
  nonExistingRepository=git@github-testing:janisZisenisTesting/not-existing.git
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
        \"$nonExistingRepository\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_output "Repository $nonExistingRepository was not found. Skipping it"
}

@test "if team json contains non-existing and existing repository, 'setup' should clone the existing one" {
  repositoriesPath=repositories
  repositoryName=1_TestRepository
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/non-exising.git\",
          \"git@github-testing:janisZisenisTesting/$repositoryName.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$repositoriesPath/$repositoryName/.git"
}

@test "if team json contains repositories but running without 'setup' does not clone the repositories" {
  repositoriesPath=repositories
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/1_TestRepository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt

  assert_directory_does_not_exist "$(testEnvDir)/$repositoriesPath"
}

@test "if team json contains repositoriesPrefixes 'setup' should ignore the prefixes while cloning the repositories" {
  repositoriesPath=repositories
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPrefixes\": [
        \"Prefix1_\",
        \"Prefix2_\"
      ],
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/Prefix1_TestRepository1.git\",
          \"git@github-testing:janisZisenisTesting/Prefix2_TestRepository2.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$repositoriesPath/TestRepository1/.git"
  assert_directory_exists "$(testEnvDir)/$repositoriesPath/TestRepository2/.git"
}

@test "if team json contains repositoriesPrefixes 'setup' should not ignore the prefixes while cloning when the prefixes are not in the beginning of the repository names" {
  repositoriesPath=repositories
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPrefixes\": [
        \"TestRepository1\",
        \"TestRepository2\"
      ],
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/Prefix1_TestRepository1.git\",
          \"git@github-testing:janisZisenisTesting/Prefix2_TestRepository2.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup

  assert_directory_exists "$(testEnvDir)/$repositoriesPath/Prefix1_TestRepository1/.git"
  assert_directory_exists "$(testEnvDir)/$repositoriesPath/Prefix2_TestRepository2/.git"
}

