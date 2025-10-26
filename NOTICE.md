# NOTICE

This project utilizes several third-party dependencies. Below is a list of these dependencies along with their respective licenses.

---

# Runtime Dependencies

These are Go modules that are part of the compiled application.

## Direct Dependencies

- **github.com/fatih/color v1.18.0**
  - License: MIT
  - URL: https://github.com/fatih/color

- **github.com/spf13/cobra v1.9.1**
  - License: Apache License 2.0
  - URL: https://github.com/spf13/cobra

- **github.com/spf13/viper v1.20.1**
  - License: MIT
  - URL: https://github.com/spf13/viper

## Indirect Dependencies

- **github.com/fsnotify/fsnotify v1.9.0**
  - License: BSD-3-Clause
  - URL: https://github.com/fsnotify/fsnotify

- **github.com/go-viper/mapstructure/v2 v2.4.0**
  - License: MIT
  - URL: https://github.com/go-viper/mapstructure

- **github.com/inconshreveable/mousetrap v1.1.0**
  - License: Apache License 2.0
  - URL: https://github.com/inconshreveable/mousetrap

- **github.com/mattn/go-colorable v0.1.14**
  - License: MIT
  - URL: https://github.com/mattn/go-colorable

- **github.com/mattn/go-isatty v0.0.20**
  - License: MIT
  - URL: https://github.com/mattn/go-isatty

- **github.com/pelletier/go-toml/v2 v2.2.4**
  - License: MIT
  - URL: https://github.com/pelletier/go-toml

- **github.com/rogpeppe/go-internal v1.14.1**
  - License: BSD-3-Clause
  - URL: https://github.com/rogpeppe/go-internal

- **github.com/sagikazarmark/locafero v0.10.0**
  - License: MIT
  - URL: https://github.com/sagikazarmark/locafero

- **github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8**
  - License: MIT
  - URL: https://github.com/sourcegraph/conc

- **github.com/spf13/afero v1.14.0**
  - License: Apache License 2.0
  - URL: https://github.com/spf13/afero

- **github.com/spf13/cast v1.9.2**
  - License: MIT
  - URL: https://github.com/spf13/cast

- **github.com/spf13/pflag v1.0.7**
  - License: BSD-3-Clause
  - URL: https://github.com/spf13/pflag

- **github.com/subosito/gotenv v1.6.0**
  - License: MIT
  - URL: https://github.com/subosito/gotenv

- **golang.org/x/sys v0.35.0**
  - License: BSD-3-Clause
  - URL: https://github.com/golang/sys

- **golang.org/x/text v0.28.0**
  - License: BSD-3-Clause
  - URL: https://github.com/golang/text

- **gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c**
  - License: BSD-2-Clause
  - URL: https://github.com/go-check/check

- **gopkg.in/yaml.v3 v3.0.1**
  - License: MIT
  - URL: https://github.com/go-yaml/yaml

---

# Development and Testing Dependencies

These tools are used during development, testing, and CI/CD processes, but are not bundled with the compiled application.

## Build and Code Quality Tools

- **golangci-lint v2.1.6**
  - License: GPL-3.0
  - URL: https://github.com/golangci/golangci-lint
  - Description: Go linter with 80+ enabled linters for code quality

- **gofumpt v0.9.1**
  - License: BSD-3-Clause
  - URL: https://github.com/mvdan/gofumpt
  - Description: Stricter Go formatter (gofmt + extra rules)

- **goimports v0.38.0**
  - License: BSD-3-Clause
  - URL: https://github.com/golang/tools
  - Description: Go import formatter and organizer

- **ShellCheck v0.11.0**
  - License: GPL-3.0
  - URL: https://github.com/koalaman/shellcheck
  - Description: Bash/shell script linter

- **shfmt v3.12.0**
  - License: MIT
  - URL: https://github.com/mvdan/sh
  - Description: Bash/shell script formatter

- **jq**
  - License: MIT
  - URL: https://stedolan.github.io/jq/
  - Description: JSON command-line processor

- **bats (Bash Automated Testing System)**
  - License: MIT
  - URL: https://github.com/bats-core/bats-core
  - Description: Bash testing framework for end-to-end tests

- **bats-support v0.3.0**
    - License: MIT
    - URL: https://github.com/bats-core/bats-support
    - Description: Common test helper library

- **bats-file v0.4.0**
    - License: MIT
    - URL: https://github.com/bats-core/bats-file
    - Description: File and path assertion helpers

- **bats-assert v2.1.0**
    - License: MIT
    - URL: https://github.com/bats-core/bats-assert
    - Description: Assertion helpers for BATS tests

- **bats-detik v1.3.1**
    - License: MIT
    - URL: https://github.com/bats-core/bats-detik
    - Description: Docker and container testing helpers

- **GNU parallel**
  - License: GPL-3.0
  - URL: https://www.gnu.org/software/parallel/
  - Citation: O. Tange (2018): GNU Parallel 2018, March 2018, https://doi.org/10.5281/zenodo.1146014
  - Description: Parallel job execution for running tests in parallel

## Container and Infrastructure Tools

- **Docker / Docker Compose**
  - License: Apache License 2.0
  - URL: https://github.com/moby/moby
  - Description: Container runtime and orchestration for isolated testing environments

- **tini v0.19.0**
  - License: MIT
  - URL: https://github.com/krallin/tini
  - Description: Lightweight init system for Docker containers

## CI/CD and Automation Tools

- **GitHub Actions**
  - License: Proprietary (GitHub)
  - URL: https://github.com/features/actions
  - Description: CI/CD platform for automated testing, building, and release management

- **GitHub CLI (gh)**
  - License: MIT
  - URL: https://github.com/cli/cli
  - Description: Command-line interface for repository management and automation

---

## License Compliance Note

**Project License**: This project is licensed under the MIT License.

**Runtime vs. Development Dependencies**: The runtime dependencies (Go modules) are all compatible with the MIT license. The development and testing tools section includes some GPL-3.0 licensed software (golangci-lint, ShellCheck, and GNU parallel). These tools are used exclusively during development and CI/CD processes and are **not bundled with or compiled into the distributed binary**. They run as separate processes during testing and code quality checks, and therefore do not create a licensing conflict with the MIT license.

**Distribution**: When distributing this project as a compiled binary or package, only the runtime dependencies are included. Users installing the binary do not need to install or comply with the GPL-3.0 tools.

---

This project is a comprehensive collection of utilities and components bringing together various pieces of functionality under a common framework for ease of use and improved stability. It adheres to the licenses specified above and ensures compliance with their respective terms.
