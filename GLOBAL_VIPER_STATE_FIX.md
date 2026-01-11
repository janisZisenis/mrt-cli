# Global Viper State - Future Fix Strategy

**Status:** Accepted limitation for now (uses global mutex workaround)
**File:** `app/commands/run/runscript/command.go`
**Date:** 2026-01-11

## Current State

The current implementation uses a global `sync.Mutex` to protect Viper's global state:

```go
//nolint:gochecknoglobals
var configMutex sync.Mutex

func LoadCommandConfig(commandPath string) CommandConfig {
	configMutex.Lock()
	defer configMutex.Unlock()
	// ... Viper operations
}
```

This works but is a workaround, not a proper solution. The linters rightfully complained about:
- Global variable anti-pattern
- `os.Exit()` preventing defer cleanup

## The Root Problem

**Viper uses global internal state.** When you call:

```go
viper.SetConfigName("config")
viper.SetConfigType("json")
viper.AddConfigPath("/some/path")
viper.ReadInConfig()
```

These operations modify **global state inside Viper**, not creating isolated config instances. This means:

1. Multiple concurrent calls interfere with each other
2. No way to safely have multiple configs loaded simultaneously
3. Forces global synchronization (mutexes) as a workaround

## Proper Fix Strategy: Use Separate Viper Instances

Instead of fighting Viper's design, work **with** it by creating separate instances:

### Step 1: Understand Viper Instances

Viper allows creating isolated instances without global state:

```go
// Instead of global Viper, create instances
v := viper.New()
v.SetConfigName("config")
v.SetConfigType("json")
v.AddConfigPath(commandDir)

if err := v.ReadInConfig(); err != nil {
    // Handle error
}

var config CommandConfig
if err := v.Unmarshal(&config); err != nil {
    // Handle error
}
```

### Step 2: Refactor LoadCommandConfig

**Before (current):**
```go
var configMutex sync.Mutex  // Global mutex needed!

func LoadCommandConfig(commandPath string) CommandConfig {
	configMutex.Lock()
	defer configMutex.Unlock()

	commandDir := filepath.Dir(commandPath)
	setupDefaults(commandDir)

	viper.AddConfigPath(commandDir)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	readErr := viper.ReadInConfig()
	if readErr != nil {
		if errors.As(readErr, &viper.ConfigFileNotFoundError{}) {
			return defaultConfig(commandDir)
		}
		log.Errorf("Error while reading config: %v", readErr)
		os.Exit(1)
	}

	var config CommandConfig
	_ = viper.Unmarshal(&config)
	return config
}
```

**After (proposed):**
```go
// No global mutex needed!

func LoadCommandConfig(commandPath string) (CommandConfig, error) {
	commandDir := filepath.Dir(commandPath)

	// Create isolated Viper instance
	v := viper.New()
	v.SetDefault("shortDescription", defaultConfig(commandDir).ShortDescription)

	// Configure with safe defaults
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(commandDir)

	// Read config
	err := v.ReadInConfig()
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return defaultConfig(commandDir), nil
		}
		return CommandConfig{}, fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal into struct
	var config CommandConfig
	if err := v.Unmarshal(&config); err != nil {
		return CommandConfig{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}
```

### Step 3: Benefits of This Approach

✅ **No global state** - Each call gets its own Viper instance
✅ **Truly concurrent-safe** - No mutexes needed
✅ **Proper error handling** - Returns errors instead of `os.Exit()`
✅ **Testable** - Can inject different configs for testing
✅ **No nolint hacks** - Code passes linting cleanly
✅ **Follows Go idioms** - Error handling, not side effects

### Step 4: Cascading Changes

This change would require updates to callers:

**Before:**
```go
config := LoadCommandConfig(scriptPath)  // Panics on error
```

**After:**
```go
config, err := LoadCommandConfig(scriptPath)
if err != nil {
	log.Errorf("Failed to load config: %v", err)
	return err  // Let caller handle
}
```

## Implementation Roadmap

### Phase 1: Refactor LoadCommandConfig (30 mins)
- Create new version using `viper.New()`
- Return `(CommandConfig, error)` instead of just `CommandConfig`
- Write unit tests for error cases

### Phase 2: Update Callers (20 mins)
- Find all calls to `LoadCommandConfig()`
- Update to handle error return
- Remove `os.Exit()` calls

### Phase 3: Remove Workaround (5 mins)
- Delete `var configMutex sync.Mutex`
- Remove nolint comments
- Verify linting passes

### Phase 4: Testing (15 mins)
- Run `mrt run lint-go` - should pass cleanly
- Test with `go run -race ./app` - no race detector warnings
- Add unit tests for concurrent config loading

**Total estimated time: ~70 minutes**

## Why Wait?

This fix is marked as "accept limitation for now" because:

1. **Current code works** - The mutex solution prevents actual races
2. **Low risk to change later** - Isolated refactor with clear test cases
3. **Allows focusing on other issues** - Error handling fixes are higher priority
4. **Can be done incrementally** - No emergency pressure

When you're ready to tackle this (maybe Phase 3 or 4 of your sprint), this roadmap provides a clear path forward.

## References

- **Viper Package:** https://pkg.go.dev/github.com/spf13/viper
- **Viper New():** Creates isolated instance without global state
- **Go Error Handling:** Return errors, don't call `os.Exit()`
- **Race Detector:** `go run -race ./app` validates concurrent safety

---

**Author's Note:** This document represents the "right way" to fix Viper state management. The current mutex approach is pragmatic but temporary. When resources allow, this refactor will make the code cleaner, faster, and more maintainable.