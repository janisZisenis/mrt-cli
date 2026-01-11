# Quick Reference: Remaining Issues

Fast access to unfixed critical issues.

## 🔴 CRITICAL #1: Unhandled Config Errors

**Files:** Multiple locations

**Pattern:**
```go
// BEFORE - DON'T DO THIS
teamInfo, _ := core.LoadTeamConfiguration()
scripts, _ := filepath.Glob(path)

// AFTER - DO THIS
teamInfo, err := core.LoadTeamConfiguration()
if err != nil {
    log.Errorf("Failed to load configuration: %v", err)
    return err
}

scripts, err := filepath.Glob(path)
if err != nil {
    log.Errorf("Failed to find scripts: %v", err)
    return err
}
```

---

## 🟠 SIGNIFICANT #1: Ignored Path Errors

**File:** `app/core/location.go:19, 24, 29`

**Before:**
```go
func GetAbsoluteExecutionPath() string {
    absolute, _ := filepath.Abs(GetExecutionPath())
    return absolute
}
```

**After:**
```go
func GetAbsoluteExecutionPath() (string, error) {
    return filepath.Abs(GetExecutionPath())
}
```

---

## 🟡 MINOR #12: Unmarshal Error Ignored

**File:** `app/commands/run/runscript/command.go:47`

**Before:**
```go
_ = viper.Unmarshal(&config)
```

**After:**
```go
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

### Run security scanner
```bash
gosec ./app/...
```

### Find all error suppressions
```bash
grep -r ", _" app/ --include="*.go"
grep -r "os.Exit" app/ --include="*.go"
```

---

## Priority Checklist

### Remaining Issues to Fix
- [ ] #1 - Config errors (CRITICAL) - ⏳ TODO
- [ ] #7 - Path errors (SIGNIFICANT) - ⏳ TODO
- [ ] #12 - Unmarshal error (MINOR) - ⏳ TODO
- [x] #14 - Hardcoded paths (MINOR) - ✅ FIXED

---

See `SECURITY_AND_BUGS_ANALYSIS.md` for detailed analysis and code examples.