package tests_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/outputs"
)

func Test_IfRunCommandConfigContainsShortDescription_Help_ShouldDisplayShortDescription(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	shortDescription := "A command that outputs some-output"
	f.RunFixture.WriteDummyCommand(commandName)
	f.RunFixture.WriteConfig(commandName, fixtures.WithShortDescription(shortDescription))

	output, _ := f.MakeMrtCommand().Run("-h").Execute()

	output.AssertHasLine(t, "  "+commandName+" "+shortDescription)
}

func Test_IfRunCommandConfigDoesNotContainShortDescription_Help_ShouldDisplayDefaultDescription(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	f.RunFixture.WriteDummyCommand(commandName)
	f.RunFixture.WriteConfig(commandName)

	output, _ := f.MakeMrtCommand().Run("-h").Execute()

	output.AssertHasLine(t, "  "+commandName+" Executes run command "+commandName)
}

func Test_IfRunCommandConfigIsAnEmptyFile_Help_ShouldExitWithErrorAndPrintErrorMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	f.RunFixture.WriteDummyCommand(commandName)
	f.RunFixture.WriteCorruptConfig(commandName)

	output, exitCode := f.MakeMrtCommand().Run("-h").Execute()

	assert.Equal(t, 1, exitCode)
	output.AssertInOrder(t,
		outputs.HasLine("Error while reading "+f.RunFixture.ConfigFilePath(commandName)),
		outputs.HasLineContaining("unexpected end of JSON input"),
	)
}
