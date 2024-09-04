testEnvDir() {
  echo "./testEnv"
}

setup() {
  load 'test_helper/writeTeamFile'
  load 'test_helper/ssh-authenticate'
  load 'test_helper/common'
  load 'test_helper/defaults'
  load 'test_helper/commits'
  load 'test_helper/assertFileExists'

  _common_setup "$(testEnvDir)"
  authenticate
}

teardown() {
  _common_teardown "$(testEnvDir)"
  revoke-authentication
}

@test "if additional pre-commit scripts exist 'committing' will execute them" {
  repository=1_TestRepository
  writeTeamFile "$(testEnvDir)" "{
      \"repositories\": [
          \"git@github-testing:janisZisenisTesting/$repository.git\"
      ]
  }"
  "$(testEnvDir)"/mrt setup
  additionalScriptsPath="$(testEnvDir)/$(default_repositories_dir)/$repository/hook-scripts/pre-commit"
  mkdir -p "$additionalScriptsPath"
  firstSpyFile="script1Executed"
  echo "
  #!/bin/bash
  touch $firstSpyFile
  " > "$additionalScriptsPath/script1"
  chmod +x "$additionalScriptsPath/script1"
  secondSpyFile="script2Executed"
  echo "
  #!/bin/bash
  touch $secondSpyFile
  " > "$additionalScriptsPath/script2"
  chmod +x "$additionalScriptsPath/script2"

  commit_changes "$(testEnvDir)/$(default_repositories_dir)/$repository" "no-prefix-branch" "no-prefix-message"

  assert_file_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/$firstSpyFile"
  assert_file_exists "$(testEnvDir)/$(default_repositories_dir)/$repository/$secondSpyFile"
}

