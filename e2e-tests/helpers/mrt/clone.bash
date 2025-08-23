clone_repositories_using_mrt() {
  bats_load_library 'testRepositories.bash'

	readarray -t repositoriesUrls < <(getRepositoryUrls "$@")
	clone_repository_urls_using_mrt "${repositoriesUrls[@]}"
}

clone_repository_urls_using_mrt() {
  bats_load_library 'writeTeamFile.bash'
	writeRepositoriesUrls "$@"

  bats_load_library 'mrt/setup.bash'
	mrtSetupClone
}