get_testing_repository_url() {
	local repository_name="$1"

	echo "git@github-testing:janisZisenisTesting/$repository_name.git"
}

get_repository_urls() {
	local repository_names=("$@")

	for r in "${repository_names[@]}"; do
		get_testing_repository_url "$r"
	done
}