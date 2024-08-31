#!/bin/bash -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
ProgName=$(basename "$0")

BINARY_PATH="$SCRIPT_DIR"/build/mrt

get_all_functions() {
  all_functions=$(declare -F)
  replaced="${all_functions//declare -f /}"
  replaced="${replaced//get_all_functions/}"
  replaced="${replaced//help/}"

  echo "${replaced[@]}"
}

help(){
  printf "Usage: %s <subcommand> [options]\n\n" "$ProgName"
  printf "Subcommands:\n"

  for command in $(get_all_functions)
  do
    printf "\t%s\n" "$command"
  done

  printf "\n"
  printf "For help with each subcommand run:\n"
  printf "%s <subcommand> -h|--help\n" "$ProgName"
  printf "\n"
}

build(){
  cd app
  go build -v -o "$BINARY_PATH"
}

run-build(){
  build
  "$SCRIPT_DIR"/build/mrt
}

run-e2e-tests() {
  $(build)
  ./e2e-test/bats/bin/bats e2e-test/*.bats
}

execute() {
  $(build)
  ./build/mrt
}

test() {
  arr=( "" )
  arr2=( "1" "3" )

  assert_equals "${arr[*]}" "${arr2[*]}"
}

assert_equals() {
  actual=$1
  expected=$2

  if [ "$expected" != "$actual" ]
  then
    printf 'actual:\t%s\n' "$actual"
    printf 'expected:\t%s\n' "$expected"
    exit 1
  fi
}

subcommand=$1
case $subcommand in
    "" | "-h" | "--help")
        help
        ;;
    *)
        shift

        if [[ $(get_all_functions) =~ $subcommand ]]
        then
          "${subcommand}" "$@"
          exit 0
        fi

        echo "Error: '$subcommand' is not a known subcommand." >&2
        echo "       Run '$ProgName --help' for a list of known subcommands." >&2
        exit 1
        ;;
esac
