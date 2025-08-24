write_empty_json_if_file_does_not_exist() {
  local file_path="$1"

	if [[ ! -f "$file_path" ]]; then
		echo "{}" > "$file_path"
	fi
}

to_json_array() {
  local input_array=("$@")

  if [[ ${#input_array[@]} -eq 0 ]]; then
    echo "[]"
    return
  fi

  jq -Rn --argjson args "$(printf '%s\n' "${input_array[@]}" | jq -R . | jq -s .)" '$args'
}

to_json_string() {
  echo "\"$1\""
}

write_json_field() {
  local file_path="$1"
  local field_name="$2"
  local value="$3"

  write_empty_json_if_file_does_not_exist "$file_path"

  local temp_file="$file_path.tmp"

  jq --argjson value "$value" ". += {\"$field_name\": \$value}" "$file_path" > "$temp_file" \
    && mv "$temp_file" "$file_path"
}