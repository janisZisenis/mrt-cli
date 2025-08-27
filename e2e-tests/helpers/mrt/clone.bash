clone_repositories_using_mrt() {
	bats_load_library 'test_repositories.bash'

	local repositories_urls
	readarray -t repositories_urls < <(get_repository_urls "$@")
	clone_repository_urls_using_mrt "${repositories_urls[@]}"
}

clone_repository_urls_using_mrt() {
	bats_load_library 'write_team_file.bash'
	write_repositories_urls "$@"

	bats_load_library 'mrt/setup.bash'
	mrt_setup_clone
}
