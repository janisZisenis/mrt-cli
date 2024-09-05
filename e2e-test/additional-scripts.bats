testEnvDir() {
  echo "./testEnv"
}

repository=1_TestRepository

setup() {
  load 'test_helper/setupRepositories'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'
  load 'test_helper/defaults'
  load 'test_helper/commits'
  load 'test_helper/pushChanges'
  load 'test_helper/assertFileExists'
  load 'test_helper/writeSpyScript'

  _common_setup "$(testEnvDir)"
  authenticate

  setupRepositories "$(testEnvDir)" "$repository"
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

@test "if additional pre-commit scripts exist 'committing' will execute them" {
  additionalScriptsPath="$(testEnvDir)/$(default_repositories_dir)/$repository/hook-scripts/pre-commit"
  firstScriptName="script1"
  secondScriptName="script2"
  writeSpyScript "$additionalScriptsPath/$firstScriptName"
  writeSpyScript "$additionalScriptsPath/$secondScriptName"

  commit_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" "some-branch" "some-message"

  assert_spy_file_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/$firstScriptName"
  assert_spy_file_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/$secondScriptName"
}

@test "if additional commit-msg scripts exist 'committing' will execute them" {
  additionalScriptsPath="$(testEnvDir)/$(default_repositories_dir)/$repository/hook-scripts/commit-msg"
  firstScriptName="script1"
  secondScriptName="script2"
  writeSpyScript "$additionalScriptsPath/$firstScriptName"
  writeSpyScript "$additionalScriptsPath/$secondScriptName"

  commit_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" "some-branch" "some-message"

  assert_spy_file_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/$firstScriptName"
  assert_spy_file_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/$secondScriptName"
}

@test "if additional pre-push scripts exist 'pushing' will execute them" {
  branch_name="some-branch"
  commit_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" "$branch_name" "no-prefix-message"
  additionalScriptsPath="$(testEnvDir)/$(default_repositories_dir)/$repository/hook-scripts/pre-push"
  firstScriptName="script1"
  secondScriptName="script2"
  writeSpyScript "$additionalScriptsPath/$firstScriptName"
  writeSpyScript "$additionalScriptsPath/$secondScriptName"

  push_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" "$branch_name"

  assert_spy_file_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/$firstScriptName"
  assert_spy_file_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/$secondScriptName"
}

