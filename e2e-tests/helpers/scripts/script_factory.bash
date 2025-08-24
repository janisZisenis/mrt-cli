make_dummy_script() {
cat <<EOF
#!/bin/bash
EOF
}

make_stub_script() {
  local output="$1"
  local exitCode="$2"

cat <<EOF
#!/bin/bash

echo "$output"
exit $exitCode
EOF
}

spy_file_suffix() {
	echo "Executed"
}

make_spy_script() {
  local scriptName="$1"

cat <<EOF
#!/bin/bash

SCRIPT_DIR="\$( cd "\$( dirname "\${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "\$@" > "\$SCRIPT_DIR/$scriptName$(spy_file_suffix)"
EOF
}

make_script_requesting_input() {
cat <<EOF
#!/bin/bash

SCRIPT_DIR="\$( cd "\$( dirname "\${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
read -r input
touch "\$SCRIPT_DIR/\$input"
EOF
}

make_std_err_script() {
  local errorMessage="$1"
cat <<EOF
#!/bin/bash

echo "$errorMessage" 1>&2
EOF
}