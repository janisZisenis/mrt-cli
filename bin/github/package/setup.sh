#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# shellcheck source=./detection-functions.sh
. "$SCRIPT_DIR/detection-functions.sh"

os=$(detect_os)
arch=$(detect_arch)

cp "$SCRIPT_DIR/mrt-$os-$arch" "$SCRIPT_DIR/../mrt"
