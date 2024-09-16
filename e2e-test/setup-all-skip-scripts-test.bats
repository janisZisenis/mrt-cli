load 'helpers/setup'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/repositoriesPath'
load 'helpers/runMrtInTestEnvironment'
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
  scriptsDir="$testEnvironmentDir/setup"
  someScript="$scriptsDir/$someScriptName/command"
  anotherScript="$scriptsDir/$anotherScriptName/command"
  writeSpyScript "$someScript"
  writeSpyScript "$anotherScript"

  run mrt setup all "--skip-$someScriptName"

  assert_spy_file_does_not_exist "$someScript"
  assert_spy_file_exists "$anotherScript"
}

@test "if two setup scripts exist setup all with skipping the second it should only run the first" {
  someScriptName="some-script"
  anotherScriptName="another-script"
  scriptsDir="$testEnvironmentDir/setup"
  someScript="$scriptsDir/$someScriptName/command"
  anotherScript="$scriptsDir/$anotherScriptName/command"
  writeSpyScript "$someScript"
  writeSpyScript "$anotherScript"

  run mrt setup all "--skip-$anotherScriptName"

  assert_spy_file_exists "$someScript"
  assert_spy_file_does_not_exist "$anotherScript"
}

@test "if one setup scripts exists setup all with skipping the script prints out skip message" {
  scriptName="some-script"
  script="$testEnvironmentDir/setup/$scriptName/command"
  writeSpyScript "$script"

  run mrt setup all "--skip-$scriptName"

  assert_output --partial "Skipping setup script: $scriptName"
}
