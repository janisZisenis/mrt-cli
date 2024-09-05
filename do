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
  cd ..
}

run-build(){
  build
  "$SCRIPT_DIR"/build/mrt "$@"
}

run-e2e-tests() {
  if [[ $# -eq 0 ]]
  then
    files=(e2e-test/*.bats)
  else
    files=($*)
  fi

  echo "Test files to be executed:"
  printf "\t%s\n" "${files[@]}"

  build
  eval "$(ssh-agent)"
  ./e2e-test/3rdParty/bats/bin/bats ${files[*]}
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
