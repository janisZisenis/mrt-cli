get_testing_repository_url() {
	local repositoryName="$1"

	echo "git@github-testing:janisZisenisTesting/$repositoryName.git"
}

get_repository_urls() {
	local repositoryNames=("$@")

	for r in "${repositoryNames[@]}"; do
		get_testing_repository_url "$r"
	done
}