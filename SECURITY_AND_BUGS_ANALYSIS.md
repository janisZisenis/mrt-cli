# Security and Bugs Analysis Report

**Analysis Date:** 2026-01-10
**Repository:** mrt-cli
**Scope:** `/app` folder (24 Go files)

---

## Executive Summary

This report documents a comprehensive analysis of the MRT CLI codebase that identified **19 issues** ranging from critical security vulnerabilities to minor performance improvements. (14 issues fixed, 1 dismissed as inapplicable, 4 remain)

### Issue Breakdown

| Severity | Count | Status | Estimate | E2E Tests |
|----------|-------|--------|----------|-----------|
| 🔴 CRITICAL | 1 | Must fix immediately | 1.5-2 hours | ✅ Yes |
| 🔴 MAJOR | 1 | Fix within days | 1.5-2 hours | ✅ Yes |
| 🟠 SIGNIFICANT | 1 | Fix within sprint | 1-1.5 hours | ⚠️ Consider |
| 🟡 MINOR | 1 | Technical debt | 0.5-1 hour | ✅ Yes |

---

## CRITICAL ISSUES (Must Fix Immediately)

### 🔴 CRITICAL #1: Unhandled Configuration Errors

**Files:**
- `app/commands/githook/command.go:33`
- `app/commands/setup/installgithooks/command.go:22`
- `app/core/scripts.go:17`

**Severity:** CRITICAL
**Type:** Security & Error Handling
**Impact:** Configuration failures silently ignored, security checks bypassed

#### Problem

```go
// githook/command.go:33
teamInfo, _ := core.LoadTeamConfiguration()  // Error silently dropped!

// If config fails to load, teamInfo is zero-value
failIfBranchIsBlocked(teamInfo)  // Empty BlockedBranches - SECURITY ISSUE!
```

Configuration loading errors are silently ignored. Program continues with empty/default configuration, which means **branch protection can be bypassed if config fails to load**.

Similar issue in `scripts.go` line 17:
```go
scripts, _ := filepath.Glob(path)  // Error silently dropped!
```

#### Impact

- 🔓 **Security vulnerability** - branch protection bypass
- 🔓 Hooks don't execute correctly
- 🔓 Silent failures

#### Fix

**app/commands/githook/command.go (lines 33-35, 58, 47, 51):**
```go
// BEFORE - Error suppressions
teamInfo, _ := core.LoadTeamConfiguration()
hookName, _ := cmd.Flags().GetString(hookNameFlag)
repositoryPath, _ := cmd.Flags().GetString(repositoryPath)
files, _ := filepath.Glob(repositoryPath + "/hook-scripts/" + hookName + "/*")

// AFTER - With error handling
teamInfo, err := core.LoadTeamConfiguration()
if err != nil {
    log.Errorf("Failed to load team configuration: %v", err)
    cmd.SilenceUsage = true
    return
}

hookName, err := cmd.Flags().GetString(hookNameFlag)
if err != nil {
    log.Errorf("Failed to get hook-name flag: %v", err)
    cmd.SilenceUsage = true
    return
}

repositoryPath, err := cmd.Flags().GetString(repositoryPath)
if err != nil {
    log.Errorf("Failed to get repository-path flag: %v", err)
    cmd.SilenceUsage = true
    return
}

// In executeAdditionalScripts function
files, err := filepath.Glob(repositoryPath + "/hook-scripts/" + hookName + "/*")
if err != nil {
    log.Errorf("Failed to find hook scripts: %v", err)
    return
}

// Replace os.Exit() with return/cmd.SilenceUsage pattern
// Line 47: if err := prefixCommitMessage(...) { ... cmd.SilenceUsage = true; return }
// Line 51: default case - cmd.SilenceUsage = true; return
```

**app/core/scripts.go (line 17):**
```go
func ForScriptInPathDo(path string, do func(scriptPath string, scriptName string)) error {
    scripts, err := filepath.Glob(path)
    if err != nil {
        return fmt.Errorf("failed to glob scripts: %w", err)
    }

    for _, script := range scripts {
        dirPath := filepath.Dir(script)
        scriptName := filepath.Base(dirPath)
        do(script, scriptName)
    }
    return nil
}
```

---

## MAJOR ISSUES (Fix Within Days)

### 🔴 MAJOR #2: Hard Exit Calls (Poor Error Handling)

**Files:**
- `app/commands/githook/prefixCommitMessage.go:47`
- `app/commands/githook/command.go:48`
- `app/core/gitBranch.go:21`

**Severity:** MAJOR
**Type:** Error Handling
**Impact:** Ungraceful shutdowns, untestable code

#### Problem

```go
if err != nil {
    log.Errorf("The given path ... does not contain a repository.")
    os.Exit(1)  // ← Hard crash! No cleanup, no error propagation
}
```

`os.Exit(1)` immediately terminates the entire process. Inconsistent with rest of codebase, prevents error handling and testing.

#### Impact

- ❌ Ungraceful shutdowns
- ❌ Deferred cleanup may not run
- ❌ Untestable code
- ❌ Caller has no chance to handle error

#### Fix

```go
// Instead of os.Exit(1), return error
func GetCurrentBranchShortName(repoDir string) (string, error) {
    var stdout bytes.Buffer
    err := NewCommandBuilder().
        WithCommand("git").
        WithArgs("-C", repoDir, "rev-parse", "--abbrev-ref", "HEAD").
        WithStdout(&stdout).
        Run()
    if err != nil {
        return "", fmt.Errorf("repository not found at %s: %w", repoDir, err)
    }

    branchName := strings.TrimSpace(stdout.String())
    if branchName == "" {
        return "", errors.New("could not determine current branch")
    }

    return branchName, nil
}
```

---

## SIGNIFICANT ISSUES (Fix Within Sprint)

### 🟠 SIGNIFICANT #1: Ignored Path Errors

**File:** `app/core/location.go:31, 36`

**Severity:** SIGNIFICANT
**Type:** Error Handling
**Impact:** Silent failures, empty paths used in subsequent operations

#### Problem

```go
func GetAbsoluteExecutionPath() string {
    absolute, _ := filepath.Abs(GetExecutionPath())  // Error ignored!
    return absolute  // Could be empty!
}

func GetExecutableName() string {
    executable, _ := os.Executable()  // Error ignored!
    return filepath.Base(executable)  // Could be empty!
}
```

Errors silently dropped. Functions can return empty strings, which fails unpredictably.

#### Fix

```go
func GetAbsoluteExecutionPath() (string, error) {
    return filepath.Abs(GetExecutionPath())
}

func GetExecutableName() (string, error) {
    executable, err := os.Executable()
    if err != nil {
        return "", fmt.Errorf("failed to get executable name: %w", err)
    }
    return filepath.Base(executable), nil
}
```

---

## MINOR ISSUES (Technical Debt)

### 🟡 MINOR #12: Unmarshal Error Ignored

**File:** `app/commands/run/runscript/command.go:47`

**Severity:** MINOR
**Type:** Error Handling

```go
_ = viper.Unmarshal(&config)
```

If config JSON is malformed, error is silently ignored.

#### Fix

```go
if err := viper.Unmarshal(&config); err != nil {
    log.Errorf("Failed to parse command config: %v", err)
    return CommandConfig{}, err
}
```

---

## Summary Table

| ID | Priority | Category | File | Issue | Status |
|----|----------|----------|------|-------|--------|
| #1 | CRITICAL | Security | githook/command.go:33 | Unhandled config errors | ⏳ TODO |
| #2 | MAJOR | Error Handling | githook/command.go:47 | Hard exit calls (os.Exit) | ⏳ TODO |
| #7 | SIGNIFICANT | Error Handling | location.go | Ignored path errors | ⏳ TODO |
| #12 | MINOR | Error Handling | runscript/command.go:47 | Unmarshal error ignored | ⏳ TODO |

---

## Recommended Fix Priority

### Phase 1: CRITICAL (Next 1-2 hours)
```
[ ] #1 - Error handling for LoadTeamConfiguration()
```

### Phase 2: MAJOR (Next 1-2 days)
```
[ ] #2 - Replace os.Exit() with error returns
```

### Phase 3: SIGNIFICANT (Within sprint)
```
[ ] #7 - Handle path errors properly
```

### Phase 4: MINOR (Technical debt)
```
[ ] #12 - Unmarshal error handling
```

---

## Testing Recommendations

### Run with race detector
```bash
go run -race ./app
```

### Run with security scanner
```bash
gosec ./app/...
```

### E2E Tests for Fixes

#### CRITICAL #1 - Unhandled Config Errors ✅ E2E Test Required
Add tests for:
- git-hook command with missing team.json configuration
- github-hook command with invalid JSON in team.json
- Scripts glob with invalid path patterns
- Expected: Error messages logged, command exits with non-zero status

#### MAJOR #2 - Hard Exit Calls ✅ E2E Test Required
Add tests for:
- git-hook with non-existent repository path
- prefixCommitMessage with missing commit message file
- Verify errors are properly reported, not hard crashes
- Expected: Error handling allows e2e test to continue, no process termination

#### SIGNIFICANT #1 - Ignored Path Errors ⚠️ Consider E2E Test
Platform-specific path operations might not need e2e tests due to OS differences.
Consider unit tests instead for mocking os.Executable() failures.

#### MINOR #12 - Unmarshal Error Ignored ✅ E2E Test Required
Add tests for:
- run command with malformed command config JSON
- setup command with invalid JSON
- Expected: Error handling gracefully, no panics

### Add unit tests for error cases
- Empty args to git hooks
- Malformed configuration
- Invalid regex patterns
- Concurrent config loading

---

## References

- OWASP Top 10: Path Traversal, Environment Variable Exposure
- Go Memory Safety: Race Conditions, Goroutine Leaks
- Go Best Practices: Error Handling, Concurrency Patterns

---

**Report Generated:** 2026-01-10
**Analysis Tool:** Claude Code Comprehensive Analysis
**Status:** 4 issues remaining (1 CRITICAL, 1 MAJOR, 1 SIGNIFICANT, 1 MINOR)