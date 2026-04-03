package tests_test

import (
	"testing"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/outputs"
	"mrt-cli/go-e2e/runcommand"
)

func Test_IfRunCommandConfigContainsShortDescription_Help_ShouldDisplayShortDescription(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	shortDescription := "A command that outputs some-output"
	w := f.RunCommandWriter()
	w.WriteDummyCommand(commandName)
	w.WriteConfig(commandName, runcommand.WithShortDescription(shortDescription))

	output, _ := f.MakeMrtCommand().Run("-h").Execute()

	output.AssertHasLine(t, "  "+commandName+" "+shortDescription)
}

func Test_IfRunCommandConfigDoesNotContainShortDescription_Help_ShouldDisplayDefaultDescription(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	w := f.RunCommandWriter()
	w.WriteDummyCommand(commandName)
	w.WriteConfig(commandName, []runcommand.ConfigOption{}...)

	output, _ := f.MakeMrtCommand().Run("-h").Execute()

	output.AssertHasLine(t, "  "+commandName+" Executes run command "+commandName)
}

func Test_IfRunCommandConfigIsAnEmptyFile_Help_ShouldExitWithErrorAndPrintErrorMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	w := f.RunCommandWriter()
	w.WriteDummyCommand(commandName)
	w.WriteCorruptConfig(commandName)

	output, exitCode := f.MakeMrtCommand().Run("-h").Execute()

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output.AssertInOrder(t,
		outputs.HasLine("Error while reading "+w.ConfigFilePath(commandName)),
		outputs.HasLineContaining("unexpected end of JSON input"),
	)
}
