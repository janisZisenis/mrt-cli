unique_branch_name() {
	echo "branch-$(uuidgen)"
}

clone_testing_repositories() {
	local repository_dir="$1"
	shift
	local repositories_to_clone=("$@")

	bats_load_library 'test_repositories.bash'
	for repository_to_clone in "${repositories_to_clone[@]}"; do
		git clone "$(get_testing_repository_url "$repository_to_clone")" "$repository_dir/$repository_to_clone"
	done
}

commit_changes() {
	local repository_path="$1"
	local branch_name="$2"
	local commit_message=${3:-"Some Commit"}

	git -C "$repository_path" checkout -b "$branch_name"
	touch "$repository_path"/some_file
	git -C "$repository_path" add .
	git -C "$repository_path" commit -m "$commit_message"
}

commit_changes_bypassing_githooks() {
	local repository_path="$1"
	local branch_name="$2"

	git -C "$repository_path" checkout -b "$branch_name"
	touch "$repository_path"/some_file
	git -C "$repository_path" add .
	git -C "$repository_path" commit -m "Some Commit" --no-verify
}

get_commit_message_of_last_commit() {
	local repository_path="$1"

	git -C "$repository_path" log -1 --pretty=%B
}

push_changes() {
	local repository_path="$1"
	local branch_name="$2"

	run git -C "$repository_path" push --set-upstream origin "$branch_name"

	if [[ $(git -C "$repository_path" ls-remote --heads origin "$branch_name") ]]; then
		git -C "$repository_path" push origin --delete "$branch_name" --no-verify
	fi
}
