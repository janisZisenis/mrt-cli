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

## Testing Commands

### Find all error suppressions
```bash
grep -r ", _" app/ --include="*.go"
grep -r "os.Exit" app/ --include="*.go"
```

---

## Priority Checklist

- [ ] SIGNIFICANT #1 - Ignored config error (`installgithooks/command.go`)