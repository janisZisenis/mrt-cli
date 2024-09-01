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

testEnvDir() {
  echo "./testEnv"
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
