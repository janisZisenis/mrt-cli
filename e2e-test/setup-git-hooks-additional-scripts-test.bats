load 'helpers/test-case-with-1-cloned-repository-and-set-up-git-hooks'
load 'helpers/fileAssertions'
load 'helpers/writeMockScript'

@test "if pre-commit scripts exist 'committing' will execute them" {
  scriptsPath="$(repositoryDir)/hook-scripts/pre-commit"
  firstScript="$scriptsPath/script1"
  secondScript="$scriptsPath/script2"
  writeSpyScript "$firstScript"
  writeSpyScript "$secondScript"

  commit_changes "$(repositoryDir)" "some-branch" "some-message"

  assert_spy_file_exists "$firstScript"
  assert_spy_file_exists "$secondScript"
}

@test "if commit-msg scripts exits with failure 'commiting' will also fail" {
  scriptPath="$(repositoryDir)/hook-scripts/commit-msg/script"
  writeStubScript "$scriptPath" "1" "some-output"

  run commit_changes "$(repositoryDir)" "some-branch" "some-message"

  assert_failure
}

@test "if commit-msg scripts has output 'commiting' will contain the same output" {
  scriptOutput="some-output"
  scriptPath="$(repositoryDir)/hook-scripts/commit-msg/script"
  writeStubScript "$scriptPath" "0" "$scriptOutput"

  run commit_changes "$(repositoryDir)" "some-branch" "some-message"

  assert_output --partial "$scriptOutput"
}

@test "if pre-commit hook gets executed, it gets passed the git parameters" {
  scriptPath="$(repositoryDir)/hook-scripts/pre-commit/script"
  writeSpyScript "$scriptPath"

  commit_changes "$(repositoryDir)" "some-branch" "some-message"

  assert_spy_file_has_content "$scriptPath" ""
}

@test "if pre-push hook gets executed, it gets passed the git parameters" {
  scriptPath="$(repositoryDir)/hook-scripts/pre-push/script"
  writeSpyScript "$scriptPath"
  branchName="$(unique_branch_name)"
  commit_changes "$(repositoryDir)" "$branchName" "some-message"

  push_changes "$(repositoryDir)" "$branchName"

  originUrl=$(git -C "$(repositoryDir)" config --get remote.origin.url)
  remoteName=$(git remote)
  assert_spy_file_has_content "$scriptPath" "$remoteName $originUrl"
}

@test "if commit-msg hook gets executed, it gets passed the git parameters" {
  scriptPath="$(repositoryDir)/hook-scripts/commit-msg/script"
  writeSpyScript "$scriptPath"

  commit_changes "$(repositoryDir)" "some-branch" "some-message"

  assert_spy_file_has_content "$scriptPath" ".git/COMMIT_EDITMSG"
}

