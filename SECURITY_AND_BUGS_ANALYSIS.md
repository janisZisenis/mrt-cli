# Security and Bugs Analysis Report

**Analysis Date:** 2026-01-10 (last updated: 2026-04-14)
**Repository:** mrt-cli
**Scope:** `/app` folder

---

## Remaining Issues

### SIGNIFICANT #1: Ignored Config Error

**File:** `app/commands/setup/installgithooks/command.go:22`
**Type:** Error Handling

```go
// CURRENT - error silently dropped
teamInfo, _ := core.LoadTeamConfiguration(".")

// SHOULD BE
teamInfo, err := core.LoadTeamConfiguration(".")
if err != nil {
    log.Errorf("Failed to load team configuration: %v", err)
    return
}
```

---

### SIGNIFICANT #2: Ignored filepath.Glob Errors

**Files:**
- `app/commands/githook/command.go:87`
- `app/commands/setup/installgithooks/setupGitHooks.go:18`
- `app/core/scripts.go:16`

**Type:** Error Handling
**Note:** `filepath.Glob` only errors on malformed patterns, so this is low-risk but still bad practice.

```go
// CURRENT - error silently dropped
files, _ := filepath.Glob(pattern)
```

---

### MINOR #1: Unmarshal Error Ignored

**File:** `app/commands/run/runscript/command.go:56`
**Type:** Error Handling

```go
// CURRENT
_ = viper.Unmarshal(&config)

// SHOULD BE
if err := viper.Unmarshal(&config); err != nil {
    log.Errorf("Failed to parse command config: %v", err)
    return CommandConfig{}, err
}
```

---

## Fixed Issues

| ID | File | Issue | Fixed |
|----|------|-------|-------|
| CRITICAL #1 | `githook/command.go` | Unhandled config error | ✅ Already handled at line 42-46 |
| MAJOR #2 | `core/gitBranch.go` | `os.Exit` instead of returning error | ✅ Fixed 2026-04-14 |
| MAJOR #2 | `githook/prefixCommitMessage.go` | `os.Exit` instead of returning error | ✅ Already fixed |
| SIGNIFICANT #1 | `core/location.go` | Ignored path errors | ✅ File removed/refactored |