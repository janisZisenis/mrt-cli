setup() {
  load 'test_helper/assertArrayEquals'
  load 'test_helper/findRepositories'
  load 'test_helper/writeTeamFile'
  load 'test_helper/common'

  _common_setup "$(testEnvDir)"
}

teardown() {
  _common_teardown "$(testEnvDir)"
}

testEnvDir() {
  echo "./testEnv"
}

@test "if team json contains repositories 'setup --all' clones that repository into given repository path" {
  repositoriesPath=repositories
  firstRepository=BoardGames.TDD-London-School
  secondRepository=BowlingGameKata
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github.com:janisZisenis/$firstRepository.git\",
          \"git@github.com:janisZisenis/$secondRepository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup --all

  actual=( $(find_repositories "$(testEnvDir)/$repositoriesPath") )
  expected=("$(testEnvDir)/$repositoriesPath/$firstRepository/.git" "$(testEnvDir)/$repositoriesPath/$secondRepository/.git")
  assert_array_equals "${actual[*]}" "${expected[*]}"
}

@test "if team json contains xyz as repositoriesPath 'setup --all' clones the repositories into given xyz folder" {
  repositoriesPath=xyz
  repository=BoardGames.TDD-London-School
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github.com:janisZisenis/$repository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup --all

  actual=( $(find_repositories "$(testEnvDir)/$repositoriesPath") )
  expected=("$(testEnvDir)/$repositoriesPath/$repository/.git")
  assert_array_equals "${actual[*]}" "${expected[*]}"
}

@test "if team json contains already existing repositories 'setup --all' clones remaining repositories given repository path" {
  repositoriesPath=repositories
  git clone git@github.com:janisZisenis/BoardGames.TDD-London-School.git "$(testEnvDir)"/$repositoriesPath/BoardGames.TDD-London-School
  firstRepository=BoardGames.TDD-London-School
  secondRepository=BowlingGameKata
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": [
          \"git@github.com:janisZisenis/$firstRepository.git\",
          \"git@github.com:janisZisenis/$secondRepository.git\"
      ]
  }"

  run "$(testEnvDir)"/mrt setup --all

  actual=( $(find_repositories "$(testEnvDir)/$repositoriesPath") )
  expected=("$(testEnvDir)/$repositoriesPath/$firstRepository/.git" "$(testEnvDir)/$repositoriesPath/$secondRepository/.git")
  assert_array_equals "${actual[*]}" "${expected[*]}"
}

@test "if team json does not contains any repository, 'setup --all' does not clone any repository" {
  repositoriesPath=repositories
  writeTeamFile "$(testEnvDir)" "{
      \"repositoriesPath\": \"$repositoriesPath\",
      \"repositories\": []
  }"

  run "$(testEnvDir)"/mrt setup --all

  actual=( $(find_repositories "$(testEnvDir)/$repositoriesPath") )
  expected=()
  assert_array_equals "${actual[*]}" "${expected[*]}"
  assert_output "The $(teamFileName) file does not contain any repositories"
}

@test "if repository is not available, 'setup --all' should print out a message" {
  assert_array_equals "just to know where to" "start the next time"
}