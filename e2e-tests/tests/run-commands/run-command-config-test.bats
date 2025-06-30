bats_load_library 'common'
bats_load_library 'writeMockCommand'

setup() {
	common_setup
}

teardown() {
	common_teardown
}

@test "if command succeeds with output it will print the command's output" {
	commandName="some-command"
  commandLocation="$(testEnvDir)/run"
  shortDescription="A command that outputs some-output"
  writeStubCommand "$commandLocation" "$commandName" "0" "some-output"
  cat <<EOF >"$(testEnvDir)/run/$commandName/config.json"
{
  "short_description": "$shortDescription"
}
EOF

  run bats_pipe runCommand "-h" \| grep "$commandName"

	assert_output "  $commandName $shortDescription"
}