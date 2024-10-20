# Multi Repository Tool

The *Multi Repository Tool* is a command-line utility designed to streamline the management of multiple repositories within a team. It allows you to:

- Automatically clone a predefined list of team repositories.
- Execute setup scripts automatically after cloning the repositories.
- Manage Git hooks and rules across multiple or individual repositories, such as enforcing JIRA-ID based prefixes for commit messages or restricting commits and pushes to specified branches.
- Run automation scripts effortlessly on your local machine or within your CI/CD pipeline.

## Installation

Follow these steps to install the Multi-Repo Tool on your machine:

### Step 1: Download the Binary

1. Navigate to the [Releases](https://github.com/janisZisenis/multi-repo-tool/releases) section of the repository.
2. Download the appropriate binary for your platform (macOS or Linux).

### Step 2: Rename the Binary

Rename the downloaded binary to `mrt` for easier usage.

```sh
mv path/to/downloaded-binary mrt
```

### Step 3: Add to PATH

To use the `mrt` command from any location, you need to add it to your PATH.

#### macOS/Linux

1. Move the binary to a directory that's already in your PATH, or add its current directory to your PATH. For example, you can move the `mrt` binary to `/usr/local/bin`:

```sh
sudo mv mrt /usr/local/bin/
```

2. Verify that mrt is in your PATH by running:

```sh
mrt --version
``
