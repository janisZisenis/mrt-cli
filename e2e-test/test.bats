setup() {
    load 'test_helper/common-setup'
    _common_setup

    load 'test_helper/assertArrayEquals'
    load 'test_helper/findRepositories'

    mkdir -p $(testEnvDir)
    cp $(which mrt) $(testEnvDir)/mrt
}

teardown() {
    rm -rf $(testEnvDir)
}

testEnvDir() {
    echo "./testEnv"
}

@test "if team json contains repositories 'setup --all' clones that repository into given repository path" {
    repositoriesPath=repositories
    firstRepository=BoardGames.TDD-London-School
    secondRepository=BowlingGameKata
    echo "{
        \"repositoriesPath\": \"$repositoriesPath\",
        \"repositories\": [
            \"$firstRepository\",
            \"$secondRepository\"
        ]
    }" > $(testEnvDir)/team.json

    run $(testEnvDir)/mrt setup --all

    actual=$(find_repositories "$(testEnvDir)/$repositoriesPath")
    expected=("$(testEnvDir)/$repositoriesPath/$firstRepository/.git"
    "$(testEnvDir)/$repositoriesPath/$secondRepository/.git"
    )
    assert_array_equals "$actual" "$expected"
}

@test "if does not team json contains any repository, setup --all does not clone any repository" {
    repositoriesPath=repositories
    echo "{
        \"repositoriesPath\": \"$repositoriesPath\",
        \"repositories\": []
    }" > $(testEnvDir)/team.json

    run $(testEnvDir)/mrt setup --all

    actual=$(find_repositories "$(testEnvDir)/$repositoriesPath")
    expected=()
    assert_array_equals "$actual" "$expected"
    assert_output "The team.json file does not contain any repositories"
}