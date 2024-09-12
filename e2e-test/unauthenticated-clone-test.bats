load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/repositoriesPath'
load 'helpers/runSetup'
load 'helpers/writeTeamFile'
load 'helpers/runMrtInTestEnvironment'

setup() {
  _common_setup
}

teardown() {
  _common_teardown
}

@test "If team json contains repository and some repository path 'setup clone-repositories' should print out message, that it clones the repositories into given repository path" {
  test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories "some-path"
}

@test "If team json contains repository and another repository path 'setup clone-repositories' should print out message, that it clones the repositories into given repository path" {
  test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories "another-path"
}

test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories() {
  repositoryPath=$1
  repositories=(
    "1_TestRepository"
  )
  writeRepositoriesPath "$repositoryPath"

  run setupAll "${repositories[@]}"

  assert_line --index 0 "Start cloning repositories into \"$repositoryPath\""
  assert_line --index 3 "Cloning repositories done"
}

@test "If team json contains 2 repositories 'setup clone-repositories' should print out message clone done message after cloning second" {
  repositories=(
    "1_TestRepository"
    "2_TestRepository"
  )

  run setupAll "${repositories[@]}"

  assert_line --index 5 "Cloning repositories done"
}

@test "if team json contains existing repositories but authentication is missing, 'setup clone-repositories' should print message" {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"
  writeRepositoriesUrls "$repositoryUrl"

  run mrt setup clone-repositories

  assert_line --index 1 "Cloning $repositoryUrl into $(default_repositories_path)/$repository"
  assert_line --index 2 "You have no access to $repositoryUrl. Please make sure you have a valid ssh key in place."
}

@test "if team json does not contains any repository, 'setup clone-repositories' exits with error" {
  run setupAll ""

  assert_failure
  assert_output "Your team file does not contain any repositories"
}

@test "if team json does not exist, 'setup clone-repositories' exits with error" {
  run mrt setup clone-repositories

  assert_failure
  assert_output 'Could not read team file. Please make sure a "team.json" file exists next to the executable and that it follows proper JSON syntax'
}