load 'helpers/common'
load 'helpers/setupRepositories'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/commits'
load 'helpers/pushChanges'
load 'helpers/assertFileExists'
load 'helpers/writeSpyScript'

testEnvDir=$(_testEnvDir)
repository=1_TestRepository

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
  additionalScriptsPath="$testEnvDir/$(default_repositories_dir)/$repository/hook-scripts/pre-commit"
  firstScriptName="script1"
  secondScriptName="script2"
  writeSpyScript "$additionalScriptsPath/$firstScriptName"
  writeSpyScript "$additionalScriptsPath/$secondScriptName"

  commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "some-branch" "some-message"

  assert_spy_file_exists "$testEnvDir/$(default_repositories_dir)/$repository/$firstScriptName"
  assert_spy_file_exists "$testEnvDir/$(default_repositories_dir)/$repository/$secondScriptName"
}

@test "if additional commit-msg scripts exits with failure 'commiting' will also fail" {
  additionalScriptPath="$testEnvDir/$(default_repositories_dir)/$repository/hook-scripts/commit-msg/script"
  writeStubScript "$additionalScriptPath" "1" "some-output"

  run commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "some-branch" "some-message"

  assert_failure
}

@test "if additional commit-msg scripts has output 'commiting' will contain the same output" {
  scriptOutput="some-output"
  additionalScriptPath="$testEnvDir/$(default_repositories_dir)/$repository/hook-scripts/commit-msg/script"
  writeStubScript "$additionalScriptPath" "0" "$scriptOutput"

  run commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "some-branch" "some-message"

  assert_output --partial "$scriptOutput"
}

@test "if pre-commit hook gets executed, it gets passed the git parameters" {
  additionalScriptPath="$testEnvDir/$(default_repositories_dir)/$repository/hook-scripts/pre-commit/script"
  writeSpyScript "$additionalScriptPath"

  commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "some-branch" "some-message"

  assert_spy_file_has_content "$testEnvDir/$(default_repositories_dir)/$repository/script" ""
}

@test "if pre-push hook gets executed, it gets passed the git parameters" {
  additionalScriptPath="$testEnvDir/$(default_repositories_dir)/$repository/hook-scripts/pre-push/script"
  writeSpyScript "$additionalScriptPath"
  branchName="$(uuidgen)"
  commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "$branchName" "some-message"

  push_changes "$testEnvDir/$(default_repositories_dir)/$repository" "$branchName"

  originUrl=$(git -C "$testEnvDir/$(default_repositories_dir)/$repository" config --get remote.origin.url)
  remoteName=$(git remote)
  assert_spy_file_has_content "$testEnvDir/$(default_repositories_dir)/$repository/script" "$remoteName $originUrl"
}

@test "if commit-msg hook gets executed, it gets passed the git parameters" {
  additionalScriptPath="$testEnvDir/$(default_repositories_dir)/$repository/hook-scripts/commit-msg/script"
  writeSpyScript "$additionalScriptPath"

  commit_changes "$testEnvDir/$(default_repositories_dir)/$repository" "some-branch" "some-message"

  assert_spy_file_has_content "$testEnvDir/$(default_repositories_dir)/$repository/script" ".git/COMMIT_EDITMSG"
}

