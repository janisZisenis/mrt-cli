#!/usr/bin/env bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# shellcheck source=./detection-functions.sh
. "$SCRIPT_DIR/detection-functions.sh"

os=$(detect_os)
arch=$(detect_arch)

echo "detected os: $os, detected arch: $arch"
binaryLocation="$SCRIPT_DIR/mrt-$os-$arch"
targetLocation="$SCRIPT_DIR/../mrt"
if [ "$os" = "windows" ]; then
  binaryLocation+='.exe'
  targetLocation+='.exe'
fi

echo "binary location: $binaryLocation"
echo "target location: $targetLocation"

cp "$binaryLocation" "$targetLocation"
