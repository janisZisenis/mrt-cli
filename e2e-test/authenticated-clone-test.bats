load 'helpers/assertDirectoryExists'
load 'helpers/assertDirectoryDoesNotExist'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/defaults'
load 'helpers/setupRepositories'

repositoriesPath=$(default_repositories_dir)

repositoriesDir() {
  echo "$testEnvironmentDir/$repositoriesPath"
}

setup() {
  _common_setup
  authenticate
}

teardown() {
  _common_teardown
  revoke-authentication
}

@test "if team json does not contain repositoriesPath 'setup all' clones repository into 'repositories' folder" {
  repositories=("1_TestRepository")
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")

  run setupRepositories "${repositoriesUrls[@]}"

  assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains an existing repository 'setup all' should print a messages about successful cloning" {
  repositories=( "1_TestRepository")
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")

  run setupRepositories "${repositoriesUrls[@]}"

  assert_line --index 1 "Cloning ${repositoriesUrls[0]} into $repositoriesPath/${repositories[0]}"
  assert_line --index 2 "Successfully cloned ${repositoriesUrls[0]}"
}

@test "if team json contains repositoriesPath 'setup all' clones the repositories into given repositoriesPath folder" {
  repositoriesPath=xyz
  repositories=("1_TestRepository")
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")
  writeRepositoriesPath "$repositoriesPath"

  run setupRepositories "${repositoriesUrls[@]}"

  assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains already existing repositories 'setup all' clones remaining repositories and skips existing ones" {
  repositories=(
    "1_TestRepository"
    "2_TestRepository"
  )
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")
  git clone "${repositoriesUrls[0]}" "$(repositoriesDir)/${repositories[0]}"

  run setupRepositories "${repositoriesUrls[@]}"

  assert_directories_exists \
    "$(repositoriesDir)/${repositories[0]}/.git" \
    "$(repositoriesDir)/${repositories[1]}/.git"
}

@test "if team json does not contains any repository, 'setup all' does not clone any repository" {
  repositoriesUrls=()

  run setupRepositories "${repositoriesUrls[@]}"

  assert_directory_does_not_exist "$(repositoriesDir)"
}

@test "if team json contains non-existing repository, 'setup all' should print out a failure message" {
  repositories=("not-existing")
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")

  run setupRepositories "${repositoriesUrls[@]}"

  assert_line --index 1 "Cloning ${repositoriesUrls[0]} into $(default_repositories_dir)/${repositories[0]}"
  assert_line --index 2 "Repository ${repositoriesUrls[0]} was not found. Skipping it"
}

@test "if team json contains non-existing and existing repository, 'setup all' should clone the existing one" {
  repositories=("1_TestRepository" "non-existing")
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")

  run setupRepositories "${repositoriesUrls[@]}"

  assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains repositoriesPrefixes 'setup all' should ignore the prefixes while cloning the repositories" {
  repositories=(
    "Prefix1_TestRepository1"
    "Prefix2_TestRepository2"
  )
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")
  writeRepositoriesPrefixes "Prefix1_" "Prefix2_"

  run setupRepositories "${repositoriesUrls[@]}"

  assert_directories_exists \
   "$(repositoriesDir)/TestRepository1/.git" \
    "$(repositoriesDir)/TestRepository2/.git"
}

@test "if team json contains repositoriesPrefixes 'setup all' should not ignore the prefixes when the prefixes are not in the beginning of the repository names" {
  repositories=(
    "Prefix1_TestRepository1"
    "Prefix2_TestRepository2"
  )
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")
  writeRepositoriesPrefixes "TestRepository1" "TestRepository2"

  run setupRepositories "${repositoriesUrls[@]}"

  assert_directories_exists \
   "$(repositoriesDir)/${repositories[0]}/.git" \
   "$(repositoriesDir)/${repositories[1]}/.git"
}
