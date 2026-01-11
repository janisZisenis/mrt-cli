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

## 🔴 MAJOR #4: Path Traversal

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

## 🔴 MAJOR #5: Environment Variable Leakage

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

## 🟠 SIGNIFICANT #1: Glob Pattern Injection

**File:** `app/commands/githook/command.go:55`

**Before:**
```go
files, _ := filepath.Glob(repositoryPath + "/hook-scripts/" + hookName + "/*")
```

**After:**
```go
func sanitizeForGlob(input string) string {
    return filepath.Base(filepath.Clean(input))
}

files, err := filepath.Glob(
    filepath.Join(repositoryPath, "hook-scripts", sanitizeForGlob(hookName), "*"),
)
if err != nil {
    log.Errorf("Failed to find hook scripts: %v", err)
    return
}
```

---

## 🟠 SIGNIFICANT #2: Ignored Path Errors

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

## 🟠 SIGNIFICANT #3: Viper Config Race Condition

**File:** `app/commands/run/runscript/command.go:26-49`

**Before:**
```go
func LoadCommandConfig(commandPath string) CommandConfig {
    // ... no synchronization
}
```

**After:**
```go
var configMutex sync.Mutex

func LoadCommandConfig(commandPath string) (CommandConfig, error) {
    configMutex.Lock()
    defer configMutex.Unlock()
    // ... rest of function
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
- [ ] #2 - Global race (MAJOR) - ⏳ TODO
- [ ] #4 - Path traversal (MAJOR) - ⏳ TODO
- [ ] #5 - Env vars (MAJOR) - ⏳ TODO
- [ ] #6 - Glob injection (SIGNIFICANT) - ⏳ TODO
- [ ] #7 - Path errors (SIGNIFICANT) - ⏳ TODO
- [ ] #8 - Viper race (SIGNIFICANT) - ⏳ TODO
- [ ] #12 - Unmarshal error (MINOR) - ⏳ TODO
- [ ] #14 - Hardcoded paths (MINOR) - ⏳ TODO

---

See `SECURITY_AND_BUGS_ANALYSIS.md` for detailed analysis and code examples.