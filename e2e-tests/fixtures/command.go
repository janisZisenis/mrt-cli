package fixtures

import (
	"fmt"
	"mrt-cli/e2e-tests/assert"
	"os"
	"path/filepath"
	"testing"
)

const (
	commandFile    = "command"
	spyFileSuffix  = "Executed"
	dummyScript    = "#!/bin/bash\n"
	execPermission = 0o750
)

type CommandFixture struct {
	repoDir    string
	commandDir string
}

func NewCommandFixture(repoDir string, commandDir string) *CommandFixture {
	return &CommandFixture{repoDir: repoDir, commandDir: commandDir}
}

func (f *CommandFixture) CommandPath(commandName string) string {
	return filepath.Join(f.repoDir, f.commandDir, commandName, commandFile)
}

func (f *CommandFixture) spyFilePath(commandName string) string {
	return f.CommandPath(commandName) + spyFileSuffix
}

func (f *CommandFixture) writeScript(commandName string, content string) {
	path := f.CommandPath(commandName)

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		panic("command: failed to create command directory: " + err.Error())
	}

	if err := os.WriteFile(path, []byte(content), execPermission); err != nil {
		panic("command: failed to write command file: " + err.Error())
	}
}

// WriteDummyCommand writes a minimal no-op command script.
func (f *CommandFixture) WriteDummyCommand(commandName string) {
	f.writeScript(commandName, dummyScript)
}

// WriteSpyCommand writes a command script that records the arguments it receives
// into a companion file. Use AssertSpyWasCalledWith to verify the recorded args.
func (f *CommandFixture) WriteSpyCommand(commandName string) {
	script := fmt.Sprintf(`#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "$@" > "$SCRIPT_DIR/%s%s"
`, commandFile, spyFileSuffix)
	f.writeScript(commandName, script)
}

func (f *CommandFixture) AssertSpyWasCalledWith(
	t *testing.T,
	commandName string,
	expectedArgs string,
) {
	t.Helper()
	assert.FileHasContent(t, f.spyFilePath(commandName), expectedArgs+"\n")
}

func (f *CommandFixture) AssertSpyWasCalled(t *testing.T, commandName string) {
	t.Helper()
	assert.FileExists(t, f.spyFilePath(commandName))
}

func (f *CommandFixture) AssertSpyWasNotCalled(t *testing.T, commandName string) {
	t.Helper()
	assert.FileDoesNotExist(t, f.spyFilePath(commandName))
}

// WriteStubCommand writes a command script that exits with the given exit code
// and prints the given output to stdout.
func (f *CommandFixture) WriteStubCommand(commandName string, exitCode int, output string) {
	script := fmt.Sprintf("#!/bin/bash\necho %q\nexit %d\n", output, exitCode)
	f.writeScript(commandName, script)
}

// WriteStderrCommand writes a command script that prints the given message to stderr.
func (f *CommandFixture) WriteStderrCommand(commandName string, errMessage string) {
	script := fmt.Sprintf("#!/bin/bash\necho %q 1>&2\n", errMessage)
	f.writeScript(commandName, script)
}

// WriteInputCommand writes a command script that reads a line from stdin and
// creates a file named after the input in the same directory as the command.
func (f *CommandFixture) WriteInputCommand(commandName string) {
	script := `#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
read -r input
touch "$SCRIPT_DIR/$input"
`
	f.writeScript(commandName, script)
}

func (f *CommandFixture) AssertInputWasReceived(t *testing.T, commandName string, input string) {
	t.Helper()
	assert.FileExists(t, filepath.Join(filepath.Dir(f.CommandPath(commandName)), input))
}
