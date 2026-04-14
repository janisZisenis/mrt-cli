# Quick Reference: Remaining Issues

Fast access to unfixed issues.

## SIGNIFICANT #1: Ignored Config Errors

**File:** `app/commands/setup/installgithooks/command.go:22`

```go
// CURRENT - ignores error
teamInfo, _ := core.LoadTeamConfiguration(".")

// SHOULD BE
teamInfo, err := core.LoadTeamConfiguration(".")
if err != nil {
    log.Errorf("Failed to load team configuration: %v", err)
    return
}
```

---

## SIGNIFICANT #2: Ignored filepath.Glob Errors

**Files:**
- `app/commands/githook/command.go:87`
- `app/commands/setup/installgithooks/setupGitHooks.go:18`
- `app/core/scripts.go:16`

```go
// CURRENT - ignores error
files, _ := filepath.Glob(pattern)

// Note: filepath.Glob only errors on malformed patterns,
// so this is low-risk but still bad practice.
```

---

## MINOR #1: Unmarshal Error Ignored

**File:** `app/commands/run/runscript/command.go:56`

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

## Testing Commands

### Check for race conditions
```bash
go run -race ./app
```

### Find all error suppressions
```bash
grep -r ", _" app/ --include="*.go"
grep -r "os.Exit" app/ --include="*.go"
```

---

## Priority Checklist

- [ ] SIGNIFICANT #1 - Ignored config error (`installgithooks/command.go`)
- [ ] SIGNIFICANT #2 - Ignored Glob errors (3 files)
- [ ] MINOR #1 - Unmarshal error ignored (`runscript/command.go`)