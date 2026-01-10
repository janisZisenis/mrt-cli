# Quick Reference: Critical Fixes

Fast access to critical issues and their fixes.

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

## 🔴 MAJOR #1: Global Variable Race

**File:** `app/core/location.go`

**Before:**
```go
var teamDirectory *string  // No synchronization!
```

**After:**
```go
var (
    teamDirectory *string
    mu sync.RWMutex
)

func GetExecutionPath() string {
    mu.RLock()
    defer mu.RUnlock()
    // ... use teamDirectory
}
```

---

## 🔴 MAJOR #2: Path Traversal

**File:** `app/commands/setup/clonerepositories/cloneRepositories.go:23`

**Before:**
```go
func getRepositoryName(repositoryURL string) string {
    return strings.TrimSuffix(repositoryURL[strings.LastIndex(repositoryURL, "/")+1:], ".git")
}
// Doesn't validate against path traversal (../)
```

**After:**
```go
func getRepositoryName(repositoryURL string) string {
    lastSlash := strings.LastIndex(repositoryURL, "/")
    if lastSlash == -1 {
        return ""
    }
    name := repositoryURL[lastSlash+1:]
    name = strings.TrimSuffix(name, ".git")

    // Reject path traversal
    if strings.Contains(name, "..") || strings.Contains(name, "/") {
        return ""
    }
    return name
}
```

---

## 🔴 MAJOR #4: Environment Variable Leakage

**File:** `app/core/commandbuilder.go:61`

**Before:**
```go
cmd.Env = os.Environ()  // Leaks ALL credentials!
```

**After:**
```go
cmd.Env = []string{
    "PATH=" + os.Getenv("PATH"),
    "HOME=" + os.Getenv("HOME"),
    "USER=" + os.Getenv("USER"),
    "SSH_AUTH_SOCK=" + os.Getenv("SSH_AUTH_SOCK"),
    // Only safe variables - NO credentials!
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

### Remaining Critical Issues
- [ ] #1 - Config errors (CRITICAL) - ⏳ TODO
- [ ] #2 - Global race (MAJOR) - ⏳ TODO
- [x] #3 - File perms (MAJOR) - ✅ FIXED
- [ ] #4 - Path traversal (MAJOR) - ⏳ TODO
- [ ] #5 - Env vars (MAJOR) - ⏳ TODO

### Additional Fixes Completed
- [x] Array bounds check (prefixCommitMessage.go) - ✅ FIXED
- [x] MustCompile → regexp.Compile (prefixCommitMessage.go) - ✅ FIXED
- [x] Unbuffered pipe deadlock (gitClone.go with bufferSize) - ✅ FIXED
- [x] Inefficient strings (getFolderName with TrimPrefix) - ✅ FIXED
- [x] Cobra Execute error handling (main.go) - ✅ FIXED
- [x] Darwin flock installation (.github/workflows) - ✅ FIXED
- [x] Version build flag (.github/actions/build/action.yml) - ✅ FIXED
- [x] Build script typo correction - ✅ FIXED

---

See `SECURITY_AND_BUGS_ANALYSIS.md` for detailed analysis and code examples.