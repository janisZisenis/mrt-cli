make_dummy_script() {
cat <<EOF
#!/bin/bash
EOF
}

make_stub_script() {
  local output="$1"
  local exit_code="$2"

cat <<EOF
#!/bin/bash

echo "$output"
exit $exit_code
EOF
}

spy_file_suffix() {
	echo "Executed"
}

make_spy_script() {
  local script_name="$1"

cat <<EOF
#!/bin/bash

SCRIPT_DIR="\$( cd "\$( dirname "\${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "\$@" > "\$SCRIPT_DIR/$script_name$(spy_file_suffix)"
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
  local error_message="$1"
cat <<EOF
#!/bin/bash

echo "$error_message" 1>&2
EOF
}