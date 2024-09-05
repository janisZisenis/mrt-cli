testEnvDir() {
  echo "./testEnv"
}

repository=1_TestRepository

setup() {
  load 'helpers/setupRepositories'
  load 'helpers/ssh-authenticate'
  load 'helpers/common'
  load 'helpers/defaults'
  load 'helpers/commits'
  load 'helpers/pushChanges'
  load 'helpers/assertFileExists'
  load 'helpers/writeSpyScript'

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

