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

## Clone your team's repositories

To clone your team's repositories, create a folder on your machine, navigate into it, and then create a `team.json` configuration file. The configuration file should contain a list of your repositories as shown below:

```json
{
  "repositories": [
    "git@github.com:repository1.git",
    "git@github.com:repository2.git",
    "...",
  ]
}
```

Next, run the following command:

```sh
mrt setup clone-repositories
```

This command will automatically clone the specified repositories to a *repositories* folder next to your `team.json`.

If you want to change the location of the repositories, you can set the new repositories path in your team configuration file as follows:

```json
{
  "repositoriesPath": "path/to/repository"
  "repositories": [
    "..."
  ]
}
```

The specified path can be relative of absolute.

## Remove your team prefixes

Sometimes, teams in a bigger organization add prefixes to their repository names. To keep a better overview of the cloned repositories on your machine, you can remove your team's prefix(es) from the repository names while cloning by specifying them in the team configuration file.

```json
{
  "repositoriesPrefixes": [
    "team_prefix1",
    "team_prefix2",
    "..."
  ]
}
```