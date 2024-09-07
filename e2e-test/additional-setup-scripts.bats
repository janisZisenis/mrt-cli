load 'helpers/common'
load 'helpers/ssh-authenticate'
load 'helpers/writeMockScript'
load 'helpers/defaults'
load 'helpers/commits'
load 'helpers/branches'
load 'helpers/setupRepositories'

testEnvDir="$(_testEnvDir)"
repository="1_TestRepository"
repositoryDir="$testEnvDir/$(default_repositories_dir)/$repository"

setup() {
  _common_setup "$testEnvDir"
  authenticate

  setupRepositories "$testEnvDir" "$repository"
}

#teardown() {
#  _common_teardown "$testEnvDir"
#  revoke-authentication
#}

@test "if additional setup script exists 'setup' will execute it and pass the repository path as parameter" {
  additionalScriptsDir="$repositoryDir/setup"
  setupScript="$additionalScriptsDir/setup-command/command"
  writeSpyScript "$setupScript"

  "$testEnvDir"/mrt setup

  assert_spy_file_has_content "$setupScript" "$repositoryDir"
}


