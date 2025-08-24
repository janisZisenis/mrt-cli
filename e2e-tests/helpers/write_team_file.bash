_team_file_name() {
	echo "team.json"
}

_team_file_path() {
  bats_load_library 'fixtures/common_fixture.bash'
  echo "$(test_env_dir)/$(_team_file_name)"
}


_write_to_team_file() {
  bats_load_library 'json_writer.bash'
  write_json_field "$(_team_file_path)" "$1" "$2"
}

write_blocked_branches() {
  bats_load_library 'json_writer.bash'
	_write_to_team_file "blockedBranches" "$(to_json_array "$@")"
}

write_repositories_prefixes() {
  bats_load_library 'json_writer.bash'
	_write_to_team_file "repositoriesPrefixes" "$(to_json_array "$@")"
}

write_repositories_urls() {
  bats_load_library 'json_writer.bash'
	_write_to_team_file "repositories" "$(to_json_array "$@")"
}

write_repositories_path() {
  bats_load_library 'json_writer.bash'
  _write_to_team_file "repositoriesPath" "$(to_json_string "$1")"
}

write_commit_prefix_regex() {
  bats_load_library 'json_writer.bash'
  _write_to_team_file "commitPrefixRegex" "$(to_json_string "$1")"
}