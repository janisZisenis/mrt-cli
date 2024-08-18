setup() {
    load 'test_helper/common-setup'
    _common_setup

    mkdir -p $(testEnvDir)
    cp $(which mrt) $(testEnvDir)/mrt
}

teardown() {
    rm -rf $(testEnvDir)
}

testEnvDir() {
    echo "./testEnv"
}


@test "clones repository" {
    repositoryPath=repositories
    repositoryName=BoardGames.TDD-London-School
    echo "{
        \"repositoriesPath\": \"$repositoryPath\",
        \"repositoryNames\": [
            \"$repositoryName\"
        ]
    }" > $(testEnvDir)/team.json

    run mrt setup --all

    load 'test_helper/assert_directory_exists'
    assert_directory_exists "$(testEnvDir)/$repositoryPath/$repositoryName/.git"
}
