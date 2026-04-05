package fixtures

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"mrt-cli/go-e2e/assert"
)

type HookScriptFixture struct {
	repositoryDir string
}

func NewHookScriptFixture(repositoryDir string) *HookScriptFixture {
	return &HookScriptFixture{repositoryDir: repositoryDir}
}

func (f *HookScriptFixture) ScriptPath(hookName string, scriptName string) string {
	return filepath.Join(f.repositoryDir, "hook-scripts", hookName, scriptName)
}

func (f *HookScriptFixture) writeScript(path string, content string) {
	if err := os.MkdirAll(filepath.Dir(path), execPermission); err != nil {
		panic("hookscript: failed to create directory: " + err.Error())
	}

	if err := os.WriteFile(path, []byte(content), execPermission); err != nil {
		panic("hookscript: failed to write script: " + err.Error())
	}
}

func (f *HookScriptFixture) WriteSpyScript(hookName string, scriptName string) {
	path := f.ScriptPath(hookName, scriptName)
	script := fmt.Sprintf(`#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "$@" > "$SCRIPT_DIR/%s%s"
`, scriptName, spyFileSuffix)
	f.writeScript(path, script)
}

func (f *HookScriptFixture) WriteStubScript(hookName string, scriptName string, exitCode int, output string) {
	path := f.ScriptPath(hookName, scriptName)
	script := fmt.Sprintf("#!/bin/bash\necho %q\nexit %d\n", output, exitCode)
	f.writeScript(path, script)
}

func (f *HookScriptFixture) AssertWasExecuted(t *testing.T, hookName string, scriptName string) {
	t.Helper()
	path := f.ScriptPath(hookName, scriptName)
	assert.FileExists(t, path+spyFileSuffix)
}

func (f *HookScriptFixture) AssertWasExecutedWith(
	t *testing.T,
	hookName string,
	scriptName string,
	expectedArgs string,
) {
	t.Helper()
	f.AssertWasExecuted(t, hookName, scriptName)
	path := f.ScriptPath(hookName, scriptName)
	assert.FileHasContent(t, path+spyFileSuffix, expectedArgs+"\n")
}
