load 'helpers/setup'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/repositoriesPath'
load 'helpers/executeInTestEnvironment'
load 'helpers/directoryAssertions'
load 'helpers/writeMockScript'
load 'helpers/absolutePath'


setup() {
  _common_setup
}

teardown() {
  _common_teardown
}

@test "if two setup scripts exist setup all with skipping the first it should only run the second" {
  someScriptName="some-script"
  anotherScriptName="another-script"
  scriptLocation="$testEnvDir/setup"
  writeSpyScriptToLocation "$scriptLocation" "$someScriptName"
  writeSpyScriptToLocation "$scriptLocation" "$anotherScriptName"

  run execute setup all "--skip-$someScriptName"

  assert_spy_file_in_location_does_not_exist "$scriptLocation" "$someScriptName"
  assert_spy_file_in_location_exists "$scriptLocation" "$anotherScriptName"
}

@test "if two setup scripts exist setup all with skipping the second it should only run the first" {
  someScriptName="some-script"
  anotherScriptName="another-script"
  scriptLocation="$testEnvDir/setup"
  writeSpyScriptToLocation "$scriptLocation" "$someScriptName"
  writeSpyScriptToLocation "$scriptLocation" "$anotherScriptName"

  run execute setup all "--skip-$anotherScriptName"

  assert_spy_file_in_location_exists "$scriptLocation" "$someScriptName"
  assert_spy_file_in_location_does_not_exist "$scriptLocation" "$anotherScriptName"
}

@test "if one setup scripts exists setup all with skipping the script prints out skip message" {
  scriptName="some-script"
  scriptLocation="$testEnvDir/setup"
  writeSpyScriptToLocation "$scriptLocation" "$scriptName"

  run execute setup all "--skip-$scriptName"

  assert_output --partial "Skipping setup script: $scriptName"
}
