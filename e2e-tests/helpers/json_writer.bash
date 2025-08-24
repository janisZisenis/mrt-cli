write_empty_json_if_file_does_not_exist() {
  local filePath="$1"

	if [[ ! -f "$filePath" ]]; then
		echo "{}" > "$filePath"
	fi
}

to_json_array() {
  local inputArray=("$@")

  if [[ ${#inputArray[@]} -eq 0 ]]; then
    echo "[]"
    return
  fi

  jq -Rn --argjson args "$(printf '%s\n' "${inputArray[@]}" | jq -R . | jq -s .)" '$args'
}

to_json_string() {
  echo "\"$1\""
}

write_json_field() {
  local filePath="$1"
  local fieldName="$2"
  local value="$3"

  write_empty_json_if_file_does_not_exist "$filePath"

  local temp_file="$filePath.tmp"

  jq --argjson value "$value" ". += {\"$fieldName\": \$value}" "$filePath" > "$temp_file" \
    && mv "$temp_file" "$filePath"
}