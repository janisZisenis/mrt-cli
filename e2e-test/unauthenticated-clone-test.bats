load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/setupRepositories'
load 'helpers/writeTeamFile'

testEnvDir=$(_testEnvDir)

setup() {
  _common_setup "$testEnvDir"
}

teardown() {
  _common_teardown "$testEnvDir"
}

@test "If team json contains repository and some repository path 'setup all' should print out message, that it clones the repositories into given repository path" {
  test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories "some-path"
}

@test "If team json contains repository and another repository path 'setup all' should print out message, that it clones the repositories into given repository path" {
  test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories "another-path"
}

test_if_team_file_contains_repository_setup_prints_message_about_cloning_repositories() {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"
  repositoryPath=$1
  writeRepositoriesPath "$testEnvDir" "$repositoryPath"

  run setupRepositories "$testEnvDir" "$repositoryUrl"

  assert_line --index 0 "Start cloning repositories into \"$repositoryPath\""
  assert_line --index 3 "Cloning repositories done"
}

@test "If team json contains 2 repositories 'setup all' should print out message clone done message after cloning second" {
  run setupRepositories "$testEnvDir" \
    "$(getTestingRepositoryUrl "1_TestRepository")" \
    "$(getTestingRepositoryUrl "2_TestRepository")"

  assert_line --index 5 "Cloning repositories done"
}

@test "if team json contains existing repositories but authentication is missing, 'setup all' should print message" {
  repository="1_TestRepository"
  repositoryUrl="$(getTestingRepositoryUrl "$repository")"

  run setupRepositories "$testEnvDir" "$repositoryUrl"

  assert_line --index 1 "Cloning $repositoryUrl into $(default_repositories_dir)/$repository"
  assert_line --index 2 "You have no access to $repositoryUrl. Please make sure you have a valid ssh key in place."
}

@test "if team json does not contains any repository, 'setup all' exits with error" {
  run setupRepositories "$testEnvDir" ""

  assert_failure
  assert_output "Your team file does not contain any repositories"
}

@test "if team json does not exist, 'setup all' exits with error" {
  run "$testEnvDir"/mrt setup all

  assert_failure
  assert_output 'Could not read team file. Please make sure a "team.json" file exists next to the executable and that it follows proper JSON syntax'
}