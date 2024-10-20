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

### Remove your team prefixes

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

## Install git hooks

The *Multi Repository Tool* let's you also manage the git-hooks of all your cloned repositories. By running the following command you can install git-hooks to all the repositories located in your `repositoriesPath` (by default *./repositories*):

```sh
    mrt setup install-git-hooks
```

When you perform some actions in your repositories (e.g. committing/pushing) the respected git-hooks are called. These git-hooks execute the tool's `git-hook` subcommand passing their git-hook name (e.g. pre-commit/pre-push) and their repository's root path.

Currently, the *Multi Repository Tool* only supports the following git-hooks:
- commit-msg
- pre-commit
- post-commit


### Block your most important branches

Delegating the execution of the git-hook to the *Multi Repository Tool* also enables it to add automated rules across all the repositories located in your `repositoriesPath`.

Adding `blockedBranches` to your team configuration file as shown below will make commiting and pushing fail on the specified branches.

```json
{
  "repositories": [
    "..."
  ],
  "blockedBranches": [
    "main",
    "dev"
  ]
}
```

### Prefix your commit messages

Often, teams use a ticketing systems such as JIRA. To keep a good overview of which commit is implemented as part of which ticket, the *Multi Repository Tool* assists you in prefixing your commit message with the ticket number (or JIRA Issue Id). This is also a prerequisite to enable JIRA Smart Commits (see [here](https://support.atlassian.com/bitbucket-cloud/docs/use-smart-commits/)).

To ensure your commit messages are prefixed correctly you can add a regular expression to your team configuration file as shown below.

```json
{
  "repositories": [
    "..."
  ],
  "commitPrefixRegex": "ABCD-[0-9]+"
}
```

Whenever you try to commit with a message not having a prefix conforming to the regular expression followed by a colon and a space, the commit will fail.

> **Example**:<br>
> In this case a valid commit message would be "ABCD-99:&nbsp;Some Commit". An invalid message would be "Some commit".

You can add the prefix manually to the commit message, or let the *Multi Repository Tool* take care of it. To use the tool's automation add the ticket number to the branch's name. Based on the given regular expression the tool will parse the first matching substring from the branch name and use it as prefix for your commit message.

In case you have the ticket number in your branch name and provide another one manually while comitting, the manual passed ticket number has priority.

> **Exception**: Merge commits<br>
> The prefix rule does not apply to merge commits. Even with a `commitPrefixRegex` in the team configuration file merge commits are not validated and do not need to be prefixed. Merge commits are detected by commit messages starting with "Merge branch" or "Merge remote-tracking branch".

### Add custom hook scripts

In each repository you might want to add some additional tasks if you perform a certain action in the repository (e.g. linting your code before committing). The *Multi Repository Tool* allows you to add executable scripts to the predefined path *hook-scripts/\<git-hook-name\>* within the respective repository's root folder.

Below you can see an example folder structure:

```l
/repository-root-folder
  |-- .git
  |-- <your-sources>
  |-- hook-scripts
  |   |-- pre-commit
  |   |   |-- preCommitTask1
  |   |   |-- ...
  |   |-- pre-push
  |   |   |-- prePushTask1
  |   |   |-- ...
```

## Add custom setup commands

Commonly, before a new developer in the team can start to develop, they need to follow some setup steps – usually documented in a documentation tool such as Confluence. 

> Examples:
> - Installing some tools
> - downloading and installing VPN certificates
> - setup AWS configuration file
> - ...

Unfortunately, this documentation gets out of date easily complicating the setup for the new team member.

The *Multi Repository Tool* allows you to implement custom setup *commands*. Each command consists of a command file (an executable file named "command"), which is located in a the command folder. The name of the command is specified by the name of the command folder. Once you add the command folder to the predefined path *setup* in the team folder the *Multi Repository Tool* can find and execute them as part of the `mrt setup` subcommand.

Below you can see an example folder structure:

```
/team-folder
  |-- team.json
  |-- repositories
  |   |-- ...
  |-- setup
  |   |-- install-tools
  |   |   |-- command
  |   |-- setup-aws
  |   |   |-- command
```

With the folder structure above you can run

```sh
  mrt setup install-tools
```

or

```sh
  mrt setup setup-aws
```

to execute the setup commands.