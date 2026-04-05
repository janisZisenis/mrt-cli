package fixtures

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"mrt-cli/go-e2e/assert"
)

const (
	hookScriptSpySuffix = "Executed"
	hookScriptExecPerm  = 0o750
)

// HookScriptFixture writes hook scripts into a repository's hook-scripts directory
// (e.g. repository_dir/hook-scripts/pre-commit/script).
type HookScriptFixture struct {
	repositoryDir string
}

func NewHookScriptFixture(repositoryDir string) *HookScriptFixture {
	return &HookScriptFixture{repositoryDir: repositoryDir}
}

// ScriptPath returns the absolute path to a script in the given hook's directory.
func (f *HookScriptFixture) ScriptPath(hookName string, scriptName string) string {
	return filepath.Join(f.repositoryDir, "hook-scripts", hookName, scriptName)
}

func (f *HookScriptFixture) writeScript(path string, content string) {
	if err := os.MkdirAll(filepath.Dir(path), hookScriptExecPerm); err != nil {
		panic("hookscript: failed to create directory: " + err.Error())
	}

	if err := os.WriteFile(path, []byte(content), hookScriptExecPerm); err != nil {
		panic("hookscript: failed to write script: " + err.Error())
	}
}

// WriteSpyScript writes a script that records its arguments into a companion file.
// Use AssertWasExecutedWith to verify the recorded arguments.
func (f *HookScriptFixture) WriteSpyScript(hookName string, scriptName string) {
	path := f.ScriptPath(hookName, scriptName)
	script := fmt.Sprintf(`#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "$@" > "$SCRIPT_DIR/%s%s"
`, scriptName, hookScriptSpySuffix)
	f.writeScript(path, script)
}

// WriteStubScript writes a script that exits with the given exit code and prints the given output.
func (f *HookScriptFixture) WriteStubScript(
	hookName string,
	scriptName string,
	exitCode int,
	output string,
) {
	path := f.ScriptPath(hookName, scriptName)
	script := fmt.Sprintf("#!/bin/bash\necho %q\nexit %d\n", output, exitCode)
	f.writeScript(path, script)
}

// AssertWasExecuted asserts that the spy script was called.
func (f *HookScriptFixture) AssertWasExecuted(t *testing.T, hookName string, scriptName string) {
	t.Helper()
	path := f.ScriptPath(hookName, scriptName)
	assert.FileExists(t, path+hookScriptSpySuffix)
}

// AssertWasExecutedWith asserts that the spy script was called with the given arguments.
func (f *HookScriptFixture) AssertWasExecutedWith(
	t *testing.T,
	hookName string,
	scriptName string,
	expectedArgs string,
) {
	t.Helper()
	f.AssertWasExecuted(t, hookName, scriptName)
	path := f.ScriptPath(hookName, scriptName)
	assert.FileHasContent(t, path+hookScriptSpySuffix, expectedArgs+"\n")
}
