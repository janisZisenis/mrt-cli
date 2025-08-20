bats_load_library 'common'
bats_load_library 'commands/runCommandWriter'

setup() {
	common_setup
}

teardown() {
	common_teardown
}

@test "if command config contains shortDescription, it is displayed in help" {
	commandName="some-command"
  shortDescription="A command that outputs some-output"
  writeStubRunCommand "$commandName" "0" "some-output"
  cat <<EOF >"$(testEnvDir)/run/$commandName/config.json"
{
  "shortDescription": "$shortDescription"
}
EOF

  run bats_pipe runCommand "-h" \| grep "$commandName"

	assert_output "  $commandName $shortDescription"
}

@test "if command config does not contain shortDescription the default is displayed in help" {
	commandName="some-command"
  shortDescription="A command that outputs some-output"
  writeStubRunCommand "$commandName" "0" "some-output"
  cat <<EOF >"$(testEnvDir)/run/$commandName/config.json"
{}
EOF

  run bats_pipe runCommand "-h" \| grep "$commandName"

	assert_output "  $commandName Executes run command $commandName"
}

@test "if command config is completely empty, it should exit with an error" {
	commandName="some-command"
  shortDescription="A command that outputs some-output"
  configFile="$(testEnvDir)/run/$commandName/config.json"
  writeStubRunCommand "$commandName" "0" "some-output"
  touch "$configFile"

  run runCommand "-h"

  assert_equal "$status" 1
  assert_line --index 0 "Error while reading $configFile"
  assert_line --index 1 "While parsing config: unexpected end of JSON input"
}