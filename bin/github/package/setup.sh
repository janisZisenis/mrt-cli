#!/usr/bin/env bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# shellcheck source=./detection-functions.sh
. "$SCRIPT_DIR/detection-functions.sh"

os=$(detect_os)
arch=$(detect_arch)

echo "detected os: $os, detected arch: $arch"
binaryLocation="$SCRIPT_DIR/mrt-$os-$arch"

if [ "$os" = "windows" ]; then
  binaryLocation+='.exe'
fi

cp $binaryLocation "$SCRIPT_DIR/../mrt"
