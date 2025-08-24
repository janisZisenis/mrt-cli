clone_repositories_using_mrt() {
  bats_load_library 'test_repositories.bash'

	readarray -t repositoriesUrls < <(get_repository_urls "$@")
	clone_repository_urls_using_mrt "${repositoriesUrls[@]}"
}

clone_repository_urls_using_mrt() {
  bats_load_library 'write_team_file.bash'
	write_repositories_urls "$@"

  bats_load_library 'mrt/setup.bash'
	mrt_setup_clone
}