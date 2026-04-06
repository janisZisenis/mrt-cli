package tests_test

import (
	"strconv"
	"testing"

	"mrt-cli/e2e-tests/fixtures"
)

func Test_IfSetupCommandExists_ExecutingIt_WillPassTheTeamFolderAsParameter(t *testing.T) {
	tests := []string{"some-command", "another-command"}

	for _, commandName := range tests {
		t.Run(commandName, func(t *testing.T) {
			testIfSetupCommandExistsExecutingItWillPassTheTeamFolderAsParameter(t, commandName)
		})
	}
}

func testIfSetupCommandExistsExecutingItWillPassTheTeamFolderAsParameter(t *testing.T, commandName string) {
	t.Helper()
	f := fixtures.MakeMrtFixture(t).Parallel()
	f.SetupFixture.WriteSpyCommand(commandName)

	f.MakeMrtCommand().
		Setup().
		SubCommand(commandName).
		Execute()

	f.SetupFixture.AssertSpyWasCalledWith(t, commandName, f.SetupFixture.RepoDir)
}

func Test_IfSetupCommandSucceedsWithOutput_ItWillPrintTheCommandsOutput(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	someOutput := "some-output"
	f.SetupFixture.WriteStubCommand(commandName, 0, someOutput)

	output, _ := f.MakeMrtCommand().
		Setup().
		SubCommand(commandName).
		Execute()

	output.AssertLineEquals(t, 0, "Execute setup command: "+commandName)
	output.AssertLineEquals(t, 1, someOutput)
	output.AssertLineEquals(t, 2, commandName+" executed successfully")
}

func Test_IfSetupCommandFailsWithOutput_ItWillPrintTheCommandsOutputAndTheFailure(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "another-command"
	someOutput := "another-output"
	exitCode := 15
	f.SetupFixture.WriteStubCommand(commandName, exitCode, someOutput)

	output, _ := f.MakeMrtCommand().
		Setup().
		SubCommand(commandName).
		Execute()

	output.AssertLineEquals(t, 0, "Execute setup command: "+commandName)
	output.AssertLineEquals(t, 1, someOutput)
	output.AssertLineEquals(t, 2, commandName+" failed with: exit status "+strconv.Itoa(exitCode))
}

func Test_IfSetupCommandIsRequestingInput_ItShouldProcessTheInput(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "input"
	input := "some-input"
	f.SetupFixture.WriteInputCommand(commandName)

	f.MakeMrtCommand().
		Setup().
		SubCommand(commandName).
		ExecuteWithInput(input + "\n")

	f.SetupFixture.AssertInputWasReceived(t, commandName, input)
}

func Test_IfSetupCommandWritesToStderr_ItOutputsStderr(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "error"
	errMessage := "some-error"
	f.SetupFixture.WriteStderrCommand(commandName, errMessage)

	output, _ := f.MakeMrtCommand().
		Setup().
		SubCommand(commandName).
		Execute()

	output.AssertHasLine(t, errMessage)
}
