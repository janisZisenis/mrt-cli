package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/outputs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IfCommandIsRun_ItShouldPassParametersToIt(t *testing.T) {
	tests := []struct {
		name        string
		commandName string
		parameters  []string
	}{
		{"some-command with flags", "some-command", []string{"some", "--flag"}},
		{"another-command with parameters", "another-command", []string{"another", "parameter"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCommandPassesParametersToIt(t, tt.commandName, tt.parameters)
		})
	}
}

func testCommandPassesParametersToIt(
	t *testing.T,
	commandName string,
	parameters []string,
) {
	t.Helper()
	f := fixtures.MakeMrtFixture(t)
	f.RunFixture.WriteSpyCommand(commandName)

	f.MakeMrtCommandInTeamDir().
		Run().
		SubCommand(commandName, parameters...).
		Execute()

	f.RunFixture.AssertSpyWasCalled(t, commandName)
}

func Test_IfCommandSucceedsWithOutput_ItShouldPrintTheCommandsOutput(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	commandName := "some-command"
	someOutput := "some-output"
	f.RunFixture.WriteStubCommand(commandName, 0, someOutput)

	output, _ := f.MakeMrtCommandInTeamDir().
		Run().
		SubCommand(commandName).
		Execute()

	output.AssertHasLine(t, someOutput)
}

func Test_IfCommandIsRequestingInput_ItShouldProcessTheInput(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	commandName := "input"
	input := "some-input"
	f.RunFixture.WriteInputCommand(commandName)

	f.MakeMrtCommandInTeamDir().
		Run().
		SubCommand(commandName).
		ExecuteWithInput(input + "\n")

	f.RunFixture.AssertInputWasReceived(t, commandName, input)
}

func Test_IfCommandWritesToStderr_ItShouldOutputStderr(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	commandName := "error"
	errMessage := "some-error"
	f.RunFixture.WriteStderrCommand(commandName, errMessage)

	output, _ := f.MakeMrtCommandInTeamDir().
		Run().
		SubCommand(commandName).
		Execute()

	output.AssertHasLine(t, errMessage)
}

func Test_CommandExitCodeIsForwardedToTheCaller(t *testing.T) {
	tests := []struct {
		name     string
		exitCode int
	}{
		{"exits with code 0", 0},
		{"exits with code 1", 1},
		{"exits with code 2", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCommandForwardsExitCode(t, tt.exitCode)
		})
	}
}

func testCommandForwardsExitCode(t *testing.T, expectedExitCode int) {
	t.Helper()
	f := fixtures.MakeMrtFixture(t)
	commandName := "some-command"
	f.RunFixture.WriteStubCommand(commandName, expectedExitCode, "")

	_, exitCode := f.MakeMrtCommandInTeamDir().
		Run().
		SubCommand(commandName).
		Execute()

	assert.Equal(t, expectedExitCode, exitCode)
}

func Test_IfNoRunCommandsExist_RunShouldExplainHowToCreateARunCommand(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)

	output, _ := f.MakeMrtCommandInTeamDir().
		Run().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Executes a specified run command."),
		outputs.HasLine("No run commands found."),
		outputs.HasLine("To add a run command, create an executable script at:"),
		outputs.HasLine("  run/<command-name>/command"),
		outputs.HasLine("Example:"),
		outputs.HasLine("  run/build/command"),
		outputs.HasLine("  run/test/command"),
		outputs.HasLine("  run/lint/command"),
	)
}
