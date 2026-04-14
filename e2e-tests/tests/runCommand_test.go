package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
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
	f := fixtures.MakeMrtFixture(t).Parallel()
	f.RunFixture.WriteSpyCommand(commandName)

	f.MakeMrtCommandInTeamDir().
		Run().
		SubCommand(commandName, parameters...).
		Execute()

	f.RunFixture.AssertSpyWasCalled(t, commandName)
}

func Test_IfCommandSucceedsWithOutput_ItShouldPrintTheCommandsOutput(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
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
	f := fixtures.MakeMrtFixture(t).Parallel()
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
	f := fixtures.MakeMrtFixture(t).Parallel()
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
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	f.RunFixture.WriteStubCommand(commandName, expectedExitCode, "")

	_, exitCode := f.MakeMrtCommandInTeamDir().
		Run().
		SubCommand(commandName).
		Execute()

	assert.Equal(t, expectedExitCode, exitCode)
}
