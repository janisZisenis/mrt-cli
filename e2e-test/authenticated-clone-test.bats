load 'helpers/assertDirectoryExists'
load 'helpers/assertDirectoryDoesNotExist'
load 'helpers/ssh-authenticate'
load 'helpers/common'
load 'helpers/repositoriesPath'
load 'helpers/runSetup'
load 'helpers/runMrtInTestEnvironment'

repositoriesPath=$(default_repositories_path)

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

@test "if team json does not contain repositoriesPath 'setup clone-repositories' clones reposits not contain repositoriesPath 'setup clone-repositories' clones repository into 'repositories' folder" {
  repositories=("1_TestRepository")

  run setupAll "${repositories[@]}"

  assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains an existing repository 'setup clone-repositories' should print a messages about successful cloning" {
  repositories=("1_TestRepository")
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")
  writeRepositoriesUrls "${repositoriesUrls[@]}"

  run mrt setup clone-repositories

  x="[0-9]+"
  assert_line --index 1 "Cloning ${repositoriesUrls[0]} into $repositoriesPath/${repositories[0]}"
  assert_line --index 2 --regexp "Enumerating objects: $x, done."
  assert_line --index 3 --regexp "Counting objects: $x% \($x\/$x\), done."
  assert_line --index 4 --regexp "Compressing objects: $x% \($x\/$x\), done."
  assert_line --index 5 --regexp "Total [0-9]+ \(delta $x\), reused $x \(delta $x\), pack-reused $x \(from $x\)"
  assert_line --index 6 "Successfully cloned ${repositoriesUrls[0]}"
}

@test "if team json contains repositoriesPath 'setup clone-repositories' clones the repositories into given repositoriesPath folder" {
  repositoriesPath=xyz
  writeRepositoriesPath "$repositoriesPath"
  repositories=("1_TestRepository")

  run setupAll "${repositories[@]}"

  assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains already existing repositories 'setup clone-repositories' clones remaining repositories and skips existing ones" {
  repositories=(
    "1_TestRepository"
    "2_TestRepository"
  )
  git clone "$(getTestingRepositoryUrl "${repositories[0]}")" "$(repositoriesDir)/${repositories[0]}"

  run setupAll "${repositories[@]}"

  assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
  assert_directory_exists "$(repositoriesDir)/${repositories[1]}/.git"
}

@test "if team json does not contains any repository, 'setup clone-repositories' does not clone any repository" {
  run setupAll ""

  assert_directory_does_not_exist "$(repositoriesDir)"
}

@test "if team json contains non-existing repository, 'setup clone-repositories' should print out a failure message" {
  repositories=("not-existing")
  readarray -t repositoriesUrls < <(getRepositoryUrls "${repositories[@]}")
  writeRepositoriesUrls "${repositoriesUrls[@]}"

  run mrt setup clone-repositories

  assert_line --index 1 "Cloning ${repositoriesUrls[0]} into $(default_repositories_path)/${repositories[0]}"
  assert_line --index 2 "Repository ${repositoriesUrls[0]} was not found. Skipping it"
}

@test "if team json contains non-existing and existing repository, 'setup clone-repositories' should clone the existing one" {
  repositories=(
    "1_TestRepository"
    "non-existing"
  )

  run setupAll "${repositories[@]}"

  assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
}

@test "if team json contains repositoriesPrefixes 'setup clone-repositories' should ignore the prefixes while cloning the repositories" {
  repositories=(
    "Prefix1_TestRepository1"
    "Prefix2_TestRepository2"
  )
  writeRepositoriesPrefixes "Prefix1_" "Prefix2_"

  run setupAll "${repositories[@]}"

  assert_directory_exists "$(repositoriesDir)/TestRepository1/.git"
  assert_directory_exists "$(repositoriesDir)/TestRepository2/.git"
}

@test "if team json contains repositoriesPrefixes 'setup clone-repositories' should not ignore the prefixes when the prefixes are not in the beginning of the repository names" {
  repositories=(
    "Prefix1_TestRepository1"
    "Prefix2_TestRepository2"
  )
  writeRepositoriesPrefixes "TestRepository1" "TestRepository2"

  run setupAll "${repositories[@]}"

  assert_directory_exists "$(repositoriesDir)/${repositories[0]}/.git"
  assert_directory_exists "$(repositoriesDir)/${repositories[1]}/.git"
}
