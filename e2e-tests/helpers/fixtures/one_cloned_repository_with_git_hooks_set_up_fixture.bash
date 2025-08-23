bats_load_library 'fixtures/authenticated_fixture.bash'
bats_load_library 'mrt/setup.bash'
bats_load_library 'git.bash'
bats_load_library 'repositoriesPath.bash'

setup() {
	authenticated_setup

	cloneTestingRepositories "$(testEnvDir)/$(default_repositories_path)" "$(repositoryName)"
	mrtSetupGitHooks
}

teardown() {
	authenticated_teardown
}

repositoryDir() {
	echo "$(testEnvDir)/$(default_repositories_path)/$(repositoryName)"
}

repositoryName() {
  echo "1_TestRepository"
}
