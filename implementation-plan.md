# Implementation Plan: Portable Git Hooks & Team Repo Hook Installation

## Context

Currently `install-git-hooks` only installs hooks to nested/cloned repositories
(`<team-dir>/<repositoriesPath>/*/`). Hook scripts hardcode the absolute path to
the team directory, so moving the whole package requires reinstalling hooks.

This plan covers three related improvements agreed upon in design discussion:

1. **Validate `repositoriesPath`** - must be relative and must not escape the team directory
2. **Use relative paths in hook scripts** - so hooks survive moving the whole package
3. **Install hooks in the team repo itself** - if it is a git repository

---

## 1. Validate `repositoriesPath`

### Rules

- `repositoriesPath` must not be an absolute path
- `repositoriesPath` must not resolve to a location outside the team directory
  (i.e. no `../escape` patterns)

### Where

Add validation in `app/core/teamConfiguration.go` inside `LoadTeamConfiguration()`,
after the default is applied and before returning.

### Logic

```
if filepath.IsAbs(repositoriesPath) → error
if resolved path is not inside teamDir → error
```

Use `filepath.Clean` + `strings.HasPrefix` (or `filepath.Rel` and check for `..`)
to detect escape.

### Error

Return a dedicated sentinel error (e.g. `ErrInvalidRepositoriesPath`) so callers
can surface the message `"repositoriesPath must be a relative path within the team repository"`
for both invalid cases.

### Tests

New e2e tests in `e2e-tests/tests/`:

| Scenario | Exit code | Expected log output |
|---|---|---|
| `repositoriesPath` is an absolute path | non-zero | `"repositoriesPath must be a relative path within the team repository"` |
| `repositoriesPath` escapes team dir (`../outside`) | non-zero | `"repositoriesPath must be a relative path within the team repository"` |
| `repositoriesPath` is valid relative path inside team dir | 0 | succeeds as before |

Note: `ErrInvalidRepositoriesPath` is a brand new error that does not exist today.
It is only returned when `team.json` was read and unmarshalled successfully but
`repositoriesPath` fails the new validation rules. `installgithooks/command.go`
must be updated to handle this new error by exiting non-zero and logging:
`"repositoriesPath must be a relative path within the team repository"`

`ErrCouldNotReadTeamFile` (missing or corrupted `team.json`) is unaffected and
continues with defaults as before.

The default value `"repositories"` is always valid and can never trigger
`ErrInvalidRepositoriesPath` - validation only fires when `team.json` explicitly
sets a `repositoriesPath` value.

---

## 2. Use Relative Paths in Hook Scripts

### Problem

`writeGitHooks.go:19` hardcodes `getAbsoluteExecutionPath()` as `--team-dir`:

```bash
mrt --team-dir /absolute/path/to/team git-hook --hook-name "$hook_name" ...
```

Moving the package breaks all hooks.

### Solution

At install time, compute the relative path from the hook file's location back to
the team directory, and embed that in the hook script using shell path resolution:

```bash
mrt --team-dir "$(cd "$(dirname "$0")/<relative-path>" && pwd)" git-hook ...
```

The `<relative-path>` is computed once per repository during installation using
`filepath.Rel(hookFileDir, teamDir)`.

### Hook file location vs team dir

| Repository | Hook file location | Relative path to team dir |
|---|---|---|
| Nested repo (e.g. `repositoriesPath: repositories`) | `<team>/repositories/<repo>/.git/hooks/` | computed via `filepath.Rel` - depends on depth of `repositoriesPath` |
| Nested repo with `repositoriesPath: .` | `<team>/<repo>/.git/hooks/` | `../..` (same as team repo itself) |
| Team repo itself | `<team>/.git/hooks/` | always `../..` |

Because `repositoriesPath` is now guaranteed to be inside the team directory
(rule from section 1), `filepath.Rel` always produces a valid relative path.

`repositoriesPath: .` is a valid edge case - repos are direct subdirectories of
the team dir. `filepath.Rel` handles this correctly but it must be covered by tests.

### Where

Change `getHookTemplate()` in `app/commands/setup/installgithooks/writeGitHooks.go`
to accept the relative path as a parameter instead of calling
`getAbsoluteExecutionPath()` directly.

Update `writeGitHook()` to receive and pass the relative path.

Update `writeHooks()` and `setupGitHooks()` to compute the relative path per
repository and pass it down.

### Tests

Update existing hook content assertions in
`e2e-tests/tests/setupGitHooks_installsAllHooks_test.go` to verify the hook
no longer contains an absolute path. That test uses
`fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)` (defined in
`e2e-tests/fixtures/one_cloned_repository_with_git_hooks.go`) which provides
`f.ClonedRepositoryPath`.

The existing assertions check for `"#!/bin/bash -e"`, `"git-hook"`,
`"--hook-name \"$hook_name\""`, and `"--repository-path $PWD"`. After section 2,
add an assertion that the hook contains the shell expansion expression and does
not contain an absolute path:

```go
assert.FileContains(t, hookFile, `--team-dir "$(cd "$(dirname "$0")/`)
assert.FileNotContains(t, hookFile, "--team-dir /")
```

Note: `assert.FileNotContains` does not exist yet. Add it to
`e2e-tests/assert/file.go` alongside the existing `FileContains` function.

The hook resolving to the correct team dir at runtime is implicitly verified by
the existing blocked-branch and commit-prefix tests - if the relative path were
wrong, `mrt git-hook` would fail to load the team configuration.

Additional edge case tests:

| Scenario | Exit code | Expected behaviour |
|---|---|---|
| `repositoriesPath: .` | 0 | hooks work correctly; blocked branch is enforced in nested repo |

No new error paths - write failures are already covered by existing logging in
`writeGitHook`: `"Failed to create hooks directory"` and `"Failed to write hook file"`.

---

## 3. Install Hooks in the Team Repo Itself

### Rule

After installing hooks to all nested repositories, also install hooks to the team
directory itself - **if and only if `<team-dir>/.git` exists** (i.e. the team
directory is a git repository). If it is just a plain folder, skip silently.

### Where

In `app/commands/setup/installgithooks/setupGitHooks.go`, after the loop over
nested repositories, add:

```go
teamDir := getAbsoluteExecutionPath()
teamGitDir := filepath.Join(teamDir, gitMetadataDir)
if _, err := os.Stat(teamGitDir); err == nil {
    log.Infof("Installing git-hooks to team repository")
    writeHooks(teamGitDir)
    log.Successf("Done installing git-hooks to team repository.")
}
```

The relative path passed to `writeHooks` for the team repo is `../..` from
`.git/hooks/` - i.e. the team dir itself. This is simply `filepath.Rel` of
`<team>/.git/hooks` → `<team>`.

### Tests

New e2e tests:

| Scenario | Exit code | Expected log output |
|---|---|---|
| Team dir is a git repo | 0 | `"Installing git-hooks to team repository"`, `"Done installing git-hooks to team repository."` |
| Team dir is not a git repo (plain folder) | 0 | no mention of team repository in output |
| Hooks in team repo survive moving the package | 0 | hooks still execute correctly after move (use `os.Rename(f.TeamDir, newPath)` where `newPath` is a sibling directory under the same parent to avoid cross-filesystem issues on macOS; `f.TeamDir` will be stale after the move so construct git commands pointing at the new repo path directly using `f.MakeGitCommand().InDirectory(newRepoPath)`) |

Write failures during team repo hook installation are covered by the same existing
logging in `writeGitHook`: `"Failed to create hooks directory"` and `"Failed to write hook file"`.

---

## Key Context for Implementer

### `getAbsoluteExecutionPath()`

Defined in `app/commands/setup/installgithooks/environment.go`. Returns the
current working directory as an absolute path - which is always the team dir,
since `installgithooks/command.go` runs `mrt` from the team dir via
`MakeMrtCommandInTeamDir()` in tests and by convention in real usage.

### `LoadTeamConfiguration` signature and call site

```go
func LoadTeamConfiguration(teamDir string) (TeamInfo, error)
```

In `installgithooks/command.go` it is currently called as:

```go
core.LoadTeamConfiguration(".")
```

`"."` is a relative path. The validation inside `LoadTeamConfiguration` needs an
absolute team dir to resolve `repositoriesPath` against. The implementer must
convert `teamDir` to an absolute path inside `LoadTeamConfiguration` before
running validation, using `filepath.Abs(teamDir)`.

### Code snippet in section 3 is a sketch

The `writeHooks(teamGitDir)` call in the section 3 snippet is illustrative.
After section 2 is implemented, `writeHooks` will accept a relative path
parameter. For the team repo that relative path is always `../..` (from
`<team>/.git/hooks/` back to `<team>`), computed via:

```go
filepath.Rel(filepath.Join(teamDir, ".git", "hooks"), teamDir)
```

### `filepath.Rel` argument order and input values

```go
filepath.Rel(basepath, targpath)
```

- `basepath` = the hook file directory, e.g. `<team>/repositories/<repo>/.git/hooks`
- `targpath` = the team dir, e.g. `<team>`

Result is the relative path from hook location to team dir.

The glob in `setupGitHooks.go` yields paths ending in `.git` (e.g.
`<team>/repositories/<repo>/.git`). When computing `filepath.Rel`, append
`"hooks"` to the glob result to get the correct `basepath`:

```go
hookFileDir := filepath.Join(repository, "hooks") // repository is the .git path from glob
relPath, _ := filepath.Rel(hookFileDir, teamDir)
```

### Test fixture patterns

Tests use `fixtures.MakeMrtFixture(t)` which creates an isolated temp directory
as the team dir. Key methods:

```go
f := fixtures.MakeMrtFixture(t)      // creates isolated team dir in t.TempDir()
f.Authenticate()                      // adds SSH key for cloning
f.TeamConfigWriter().Write(           // writes team.json
    teamconfig.WithRepositoriesPath("some-path"),
    teamconfig.WithBlockedBranches([]string{"main"}),
)
f.AbsolutePath("some/path")          // resolves path relative to team dir
f.MakeGitCommand().Clone(url, dest)  // clones a repo
f.MakeMrtCommandInTeamDir().Setup().InstallGitHooks().Execute() // runs the command
```

Output assertions use `outputs.HasLine(...)` for exact matches and
`outputs.HasLineContaining(...)` for partial matches.

### Initialising the team dir as a git repo in tests

There is no existing fixture helper for `git init`. Use `os/exec` directly:

```go
cmd := exec.Command("git", "init", f.TeamDir)
require.NoError(t, cmd.Run())
```

### `defaultRepositoriesPath` constant

The value `"repositories"` is available in tests as a package-level const declared
in `e2e-tests/tests/authenticated_cloneRepositories_test.go:11`:

```go
const defaultRepositoriesPath = "repositories"
```

All test files in `e2e-tests/tests/` are in the same `tests_test` package so this
const is visible across all of them. Do not redeclare it in new test files.

---

## Affected Files

| File | Change |
|---|---|
| `app/core/teamConfiguration.go` | Add `repositoriesPath` validation |
| `app/commands/setup/installgithooks/command.go` | Exit non-zero on `ErrInvalidRepositoriesPath` |
| `app/commands/setup/installgithooks/writeGitHooks.go` | Accept relative team-dir path in template |
| `app/commands/setup/installgithooks/setupGitHooks.go` | Compute relative paths; install hooks in team repo |
| `e2e-tests/assert/file.go` | Add `FileNotContains` helper |
| `e2e-tests/tests/setupGitHooks_installsAllHooks_test.go` | Update hook content assertions |
| `e2e-tests/tests/setupGitHooks_defaultRepositoriesPath_test.go` | Add validation error tests |
| `e2e-tests/tests/` (new files) | Tests for team repo hook installation and relative paths |

---

## Implementation Order

1. `repositoriesPath` validation (section 1) - isolated change, unblocks the rest
2. Relative paths in hook scripts (section 2) - depends on section 1 being done first
   so the path computation can safely assume no escaping
3. Team repo hook installation (section 3) - straightforward addition once section 2 is done
