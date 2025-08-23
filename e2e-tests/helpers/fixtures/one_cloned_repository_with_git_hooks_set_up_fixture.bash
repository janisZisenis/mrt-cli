setup() {
  bats_load_library 'fixtures/authenticated_fixture.bash'
  bats_load_library 'repositoriesPath.bash'
  bats_load_library 'mrt/setup.bash'
  bats_load_library 'git.bash'

	authenticated_setup

	cloneTestingRepositories "$(testEnvDir)/$(default_repositories_path)" "$(repositoryName)"
	mrtSetupGitHooks
}

teardown() {
  bats_load_library 'fixtures/authenticated_fixture.bash'

	authenticated_teardown
}

repositoryDir() {
  bats_load_library 'repositoriesPath.bash'

	echo "$(testEnvDir)/$(default_repositories_path)/$(repositoryName)"
}

repositoryName() {
  echo "1_TestRepository"
}
