#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# shellcheck source=./detection-functions.sh
. "$SCRIPT_DIR/detection-functions.sh"

OS=$(detect_os)

sourceName="$SCRIPT_DIR/mrt-$OS-$(detect_arch)"

if [[ OS == "windows" ]]
then
    sourceName+=".exe"
fi

cp "$SCRIPT_DIR/mrt-$(detect_os)-$(detect_arch)" "$SCRIPT_DIR/../mrt"