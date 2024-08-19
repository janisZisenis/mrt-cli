#!/bin/sh

RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
NORMAL=$(tput sgr0)

expected="$1"
actual="$2"

if [ "$expected" != "$actual" ]
then
  printf "%sexpected: \t %s\n%s" "${RED}" "$expected" "${NORMAL}"
  printf "%sactual: \t %s\n%s" "${RED}" "$actual" "${NORMAL}"
  exit 1
fi