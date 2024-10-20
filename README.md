# Multi Repository Tool

The *Multi Repository Tool* is a command-line utility designed to streamline the management of multiple repositories within a team. It allows you to:

- Automatically clone a predefined list of team repositories.
- Execute setup scripts automatically after cloning the repositories.
- Manage Git hooks and rules across multiple or individual repositories, such as enforcing JIRA-ID based prefixes for commit messages or restricting commits and pushes to specified branches.
- Run automation scripts effortlessly on your local machine or within your CI/CD pipeline.

## Installation

Follow these steps to install the Multi-Repo Tool on your machine:

1. Navigate to the [Releases](https://github.com/janisZisenis/multi-repo-tool/releases) section of the repository.
2. Download the appropriate binary for your platform.
3. Make it executable.

You can achieve all of the steps above by running the script below:

```sh
#!/bin/bash

detect_os() {
    local OS
    unameOut="$(uname -s)"
    case "${unameOut}" in
        Linux*)     OS=linux;;
        Darwin*)    OS=darwin;;
        *)          OS="UNKNOWN:${unameOut}"
    esac
    echo "${OS}"
}

detect_arch() {
    local ARCH
    archOut="$(uname -m)"
    case "${archOut}" in
        x86_64)    ARCH=amd64;;
        aarch64)   ARCH=arm64;;
        arm64)     ARCH=arm64;;
        *)         ARCH="UNKNOWN:${archOut}"
    esac
    echo "${ARCH}"
}

OS=$(detect_os)
ARCH=$(detect_arch)

curl -L -o mrt https://github.com/janisZisenis/multi-repo-tool/releases/download/latest/mrt-$OS-$ARCH

chmod +x mrt
```

4. Add the binary's location to your PATH variable by running

```sh
    export PATH=$PATH:<path/to/binary>
```

5. Verify that mrt is in your PATH by running:

```sh
mrt --version

```
