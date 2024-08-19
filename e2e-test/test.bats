setup() {
    load 'test_helper/common-setup'
    _common_setup

    load 'test_helper/assertArrayEquals'
    load 'test_helper/findRepositories'

    mkdir -p "$(testEnvDir)"
    cp "$(which mrt)" "$(testEnvDir)"/mrt
}

teardown() {
    rm -rf "$(testEnvDir)"
}

testEnvDir() {
    echo "./testEnv"
}

teamFileName() {
  echo "team.json"
}

writeTeamJson() {
      echo $1 > "$(testEnvDir)/$(teamFileName)"
}

@test "if team json contains repositories 'setup --all' clones that repository into given repository path" {
    repositoriesPath=repositories
    firstRepository=BoardGames.TDD-London-School
    secondRepository=BowlingGameKata
    writeTeamJson "{
        \"repositoriesPath\": \"$repositoriesPath\",
        \"repositories\": [
            \"$firstRepository\",
            \"$secondRepository\"
        ]
    }"

    run "$(testEnvDir)"/mrt setup --all

    actual=( $(find_repositories "$(testEnvDir)/$repositoriesPath") )
    expected=("$(testEnvDir)/$repositoriesPath/$firstRepository/.git" "$(testEnvDir)/$repositoriesPath/$secondRepository/.git")
    assert_array_equals "${actual[*]}" "${expected[*]}"
}

@test "if team json contains already existing repositories 'setup --all' clones remaining repositories given repository path" {
    repositoriesPath=repositories
    git clone git@github.com:janisZisenis/BoardGames.TDD-London-School.git "$(testEnvDir)"/$repositoriesPath/BoardGames.TDD-London-School
    firstRepository=BoardGames.TDD-London-School
    secondRepository=BowlingGameKata
    writeTeamJson "{
        \"repositoriesPath\": \"$repositoriesPath\",
        \"repositories\": [
            \"$firstRepository\",
            \"$secondRepository\"
        ]
    }"

    run "$(testEnvDir)"/mrt setup --all

    actual=( $(find_repositories "$(testEnvDir)/$repositoriesPath") )
    expected=("$(testEnvDir)/$repositoriesPath/$firstRepository/.git" "$(testEnvDir)/$repositoriesPath/$secondRepository/.git")
    assert_array_equals "${actual[*]}" "${expected[*]}"
}

@test "if does not team json contains any repository, setup --all does not clone any repository" {
    repositoriesPath=repositories
    writeTeamJson "{
        \"repositoriesPath\": \"$repositoriesPath\",
        \"repositories\": []
    }"

    run "$(testEnvDir)"/mrt setup --all

    actual=( $(find_repositories "$(testEnvDir)/$repositoriesPath") )
    expected=()
    assert_array_equals "${actual[*]}" "${expected[*]}"
    assert_output "The $(teamFileName) file does not contain any repositories"
}