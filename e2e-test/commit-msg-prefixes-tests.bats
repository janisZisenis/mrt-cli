testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'
  load 'test_helper/commits'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

defaultRepositoriesPath="repositories"

@test "if team json contains jiraPrefixRegex 'commiting' with a message and a branch both without matching prefix is blocked" {
  repository=1_TestRepository
  writeTeamFile "$(testEnvDir)" "{
      \"jiraPrefixRegex\": \"Test-[0-9]+\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" "no-prefix-branch" "no-prefix-message"

  assert_line --index 1 "The commit message needs a JIRA ID prefix."
  assert_line --index 2 "Either add the JIRA ID to you commit message, or include it in the branch name."
  assert_line --index 3 "Use '--no-verify' to skip git-hooks."
  assert_failure
}

@test "if team json contains jiraPrefixRegex 'commiting' with a message that has a matching prefix on a branch without prefix is not blocked" {
  repository=1_TestRepository
  matchingPrefix="Test-1"
  writeTeamFile "$(testEnvDir)" "{
      \"jiraPrefixRegex\": \"Test-[0-9]+\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" "no-prefix-branch" "$matchingPrefix: prefixed-message"

  assert_line --index 1 "The commit message contains an issue ID ($matchingPrefix). Good job!"
  assert_success
}

@test "if team json contains jiraPrefixRegex 'commiting' with a message that has no matching prefix on a branch containing prefix is not blocked" {
  repository=1_TestRepository
  jiraId=Asdf-99
  writeTeamFile "$(testEnvDir)" "{
      \"jiraPrefixRegex\": \"Asdf-[0-9]+\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" "feature/$jiraId/prefixed-branch" "not-prefix-message"

  assert_line --index 1 "JIRA-ID '$jiraId' was found in current branch name, prepended to commit message."
  assert_success
}

@test "if team json contains jiraPrefixRegex 'commiting' with a message that starts with 'Merge branch' is not blocked" {
  test_merge_commit_messages_are_not_blocked "Merge branch"
}

@test "if team json contains jiraPrefixRegex 'commiting' with a message that starts with 'Merge remote-tracking branch' is not blocked" {
  test_merge_commit_messages_are_not_blocked "Merge remote-tracking branch"
}

@test "if team json contains jiraPrefixRegex 'commiting' with a message that has a matching prefix leads to commit with given message" {
  repository=1_TestRepository
  commitMessage="Test-1: prefixed-message"
  writeTeamFile "$(testEnvDir)" "{
      \"jiraPrefixRegex\": \"Test-[0-9]+\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup
  commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" "no-prefix-branch" "$commitMessage"

  run get_commit_message_of_last_commit "$(testEnvDir)/$defaultRepositoriesPath/$repository"
  assert_output "$commitMessage"
}

@test "if team json contains jiraPrefixRegex 'commiting' with a message that has no matching prefix on a branch containing prefix leads to commit with prefixed message" {
  repository=1_TestRepository
  matchingPrefix=Asdf-99
  commitMessage="not-prefixed-message"
  writeTeamFile "$(testEnvDir)" "{
      \"jiraPrefixRegex\": \"Asdf-[0-9]+\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup
  commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" "feature/$matchingPrefix/prefixed-branch" "$commitMessage"

  run get_commit_message_of_last_commit "$(testEnvDir)/$defaultRepositoriesPath/$repository"
  assert_output "$matchingPrefix: $commitMessage"
}

test_merge_commit_messages_are_not_blocked() {
  commit_message=$1
  repository=1_TestRepository
  writeTeamFile "$(testEnvDir)" "{
      \"jiraPrefixRegex\": \"Asdf-[0-9]+\",
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup

  run commit_changes "$(testEnvDir)/$defaultRepositoriesPath/$repository" "no-prefix-branch" "$commit_message"

  assert_line --index 1 "Merge commit detected, skipping commit-msg hook."
  assert_success
}