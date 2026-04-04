package runcommand

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"mrt-cli/go-e2e/assert"
)

const (
	runDir         = "run"
	commandFile    = "command"
	configFile     = "config.json"
	spyFileSuffix  = "Executed"
	dummyScript    = "#!/bin/bash\n"
	execPermission = 0o750
)

type Writer struct {
	teamDir string
}

func NewWriter(teamDir string) *Writer {
	return &Writer{teamDir: teamDir}
}

func (w *Writer) ConfigFilePath(commandName string) string {
	return filepath.Join(w.teamDir, runDir, commandName, configFile)
}

func (w *Writer) WriteDummyCommand(commandName string) {
	commandPath := filepath.Join(w.teamDir, runDir, commandName, commandFile)

	if err := os.MkdirAll(filepath.Dir(commandPath), 0o750); err != nil {
		panic("runcommand: failed to create command directory: " + err.Error())
	}

	if err := os.WriteFile(commandPath, []byte(dummyScript), execPermission); err != nil {
		panic("runcommand: failed to write command file: " + err.Error())
	}
}

func (w *Writer) WriteConfig(commandName string, options ...ConfigOption) {
	config := &commandConfig{}
	for _, opt := range options {
		opt(config)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic("runcommand: failed to marshal config: " + err.Error())
	}

	configPath := w.ConfigFilePath(commandName)

	if mkErr := os.MkdirAll(filepath.Dir(configPath), 0o750); mkErr != nil {
		panic("runcommand: failed to create config directory: " + mkErr.Error())
	}

	if writeErr := os.WriteFile(configPath, data, 0o600); writeErr != nil {
		panic("runcommand: failed to write config file: " + writeErr.Error())
	}
}

func (w *Writer) WriteCorruptConfig(commandName string) {
	configPath := w.ConfigFilePath(commandName)

	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		panic("runcommand: failed to create config directory: " + err.Error())
	}

	if err := os.WriteFile(configPath, []byte{}, 0o600); err != nil {
		panic("runcommand: failed to write empty config file: " + err.Error())
	}
}

func (w *Writer) commandPath(commandName string) string {
	return filepath.Join(w.teamDir, runDir, commandName, commandFile)
}

func (w *Writer) spyFilePath(commandName string) string {
	return w.commandPath(commandName) + spyFileSuffix
}

func (w *Writer) writeScript(commandName string, content string) {
	path := w.commandPath(commandName)

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		panic("runcommand: failed to create command directory: " + err.Error())
	}

	if err := os.WriteFile(path, []byte(content), execPermission); err != nil {
		panic("runcommand: failed to write command file: " + err.Error())
	}
}

// WriteSpyCommand writes a command script that records the arguments it receives
// into a companion file (commandExecuted). Use SpyFilePath to read the recorded args.
func (w *Writer) WriteSpyCommand(commandName string) {
	script := fmt.Sprintf(`#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "$@" > "$SCRIPT_DIR/%s%s"
`, commandFile, spyFileSuffix)
	w.writeScript(commandName, script)
}

func (w *Writer) AssertSpyWasCalledWith(t *testing.T, commandName string, expectedArgs string) {
	t.Helper()
	assert.FileHasContent(t, w.spyFilePath(commandName), w.teamDir+" "+expectedArgs+"\n")
}

// WriteStubCommand writes a command script that exits with the given exit code
// and prints the given output to stdout.
func (w *Writer) WriteStubCommand(commandName string, exitCode int, output string) {
	script := fmt.Sprintf("#!/bin/bash\necho %q\nexit %d\n", output, exitCode)
	w.writeScript(commandName, script)
}

// WriteStderrCommand writes a command script that prints the given message to stderr.
func (w *Writer) WriteStderrCommand(commandName string, errMessage string) {
	script := fmt.Sprintf("#!/bin/bash\necho %q 1>&2\n", errMessage)
	w.writeScript(commandName, script)
}

// WriteInputCommand writes a command script that reads a line from stdin and
// creates a file named after the input in the same directory as the command.
func (w *Writer) WriteInputCommand(commandName string) {
	script := `#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
read -r input
touch "$SCRIPT_DIR/$input"
`
	w.writeScript(commandName, script)
}

func (w *Writer) AssertInputWasReceived(t *testing.T, commandName string, input string) {
	t.Helper()
	assert.FileExists(t, filepath.Join(filepath.Dir(w.commandPath(commandName)), input))
}

type commandConfig struct {
	ShortDescription *string `json:"shortDescription,omitempty"`
}

type ConfigOption func(*commandConfig)

func WithShortDescription(description string) ConfigOption {
	return func(c *commandConfig) {
		c.ShortDescription = &description
	}
}
