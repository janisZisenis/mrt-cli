getTestingRepositoryUrl() {
	local repositoryName="$1"

	echo "git@github-testing:janisZisenisTesting/$repositoryName.git"
}

getRepositoryUrls() {
	local repositoryNames=("$@")

	for r in "${repositoryNames[@]}"; do
		getTestingRepositoryUrl "$r"
	done
}