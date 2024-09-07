load 'helpers/common'
load 'helpers/setupRepositories'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/commits'
load 'helpers/pushChanges'
load 'helpers/assertFileExists'
load 'helpers/writeMockScript'
load 'helpers/branches'

testEnvDir=$(_testEnvDir)
repository=1_TestRepository
repositoryDir="$testEnvDir/$(default_repositories_dir)/$repository"

setup() {
  _common_setup "$testEnvDir"
  authenticate

  setupRepositories "$testEnvDir" "$repository"
}

teardown() {
  _common_teardown "$testEnvDir"
  revoke-authentication
}

@test "if additional pre-commit scripts exist 'committing' will execute them" {
  additionalScriptsPath="$repositoryDir/hook-scripts/pre-commit"
  firstScriptName="$additionalScriptsPath/script1"
  secondScriptName="$additionalScriptsPath/script2"
  writeSpyScript "$firstScriptName"
  writeSpyScript "$secondScriptName"


  commit_changes "$repositoryDir" "some-branch" "some-message"

  assert_spy_file_exists "$firstScriptName"
  assert_spy_file_exists "$secondScriptName"
}

@test "if additional commit-msg scripts exits with failure 'commiting' will also fail" {
  additionalScriptPath="$repositoryDir/hook-scripts/commit-msg/script"
  writeStubScript "$additionalScriptPath" "1" "some-output"

  run commit_changes "$repositoryDir" "some-branch" "some-message"

  assert_failure
}

@test "if additional commit-msg scripts has output 'commiting' will contain the same output" {
  scriptOutput="some-output"
  additionalScriptPath="$repositoryDir/hook-scripts/commit-msg/script"
  writeStubScript "$additionalScriptPath" "0" "$scriptOutput"

  run commit_changes "$repositoryDir" "some-branch" "some-message"

  assert_output --partial "$scriptOutput"
}

@test "if pre-commit hook gets executed, it gets passed the git parameters" {
  additionalScriptPath="$repositoryDir/hook-scripts/pre-commit/script"
  writeSpyScript "$additionalScriptPath"

  commit_changes "$repositoryDir" "some-branch" "some-message"

  assert_spy_file_has_content "$additionalScriptPath" ""
}

@test "if pre-push hook gets executed, it gets passed the git parameters" {
  additionalScriptPath="$repositoryDir/hook-scripts/pre-push/script"
  writeSpyScript "$additionalScriptPath"
  branchName="$(unique_branch_name)"
  commit_changes "$repositoryDir" "$branchName" "some-message"

  push_changes "$repositoryDir" "$branchName"

  originUrl=$(git -C "$repositoryDir" config --get remote.origin.url)
  remoteName=$(git remote)
  assert_spy_file_has_content "$additionalScriptPath" "$remoteName $originUrl"
}

@test "if commit-msg hook gets executed, it gets passed the git parameters" {
  additionalScriptPath="$repositoryDir/hook-scripts/commit-msg/script"
  writeSpyScript "$additionalScriptPath"

  commit_changes "$repositoryDir" "some-branch" "some-message"

  assert_spy_file_has_content "$additionalScriptPath" ".git/COMMIT_EDITMSG"
}

