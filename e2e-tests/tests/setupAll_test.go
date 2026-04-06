package tests_test

import (
	"testing"

	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	"mrt-cli/e2e-tests/outputs"
	"mrt-cli/e2e-tests/teamconfig"
)

func Test_SetupAll_ShouldCloneInstallGitHooksAndExecuteCommands(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()
	repositoryName := "1_TestRepository"
	repositoryURL := git.MakeCloneURL(repositoryName)
	f.TeamConfigWriter().Write(teamconfig.WithRepositories([]string{repositoryURL}))
	someCommand := "some-command"
	anotherCommand := "another-command"
	f.SetupFixture.WriteSpyCommand(someCommand)
	f.SetupFixture.WriteSpyCommand(anotherCommand)

	output, _ := f.MakeMrtCommand().
		Setup().
		All().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Start cloning repositories into \"repositories\""),
		outputs.HasLine("Cloning "+repositoryURL),
		outputs.HasLineContaining("Enumerating objects:"),
		outputs.HasLine("Successfully cloned "+repositoryURL),
		outputs.HasLine("Cloning repositories done"),
		outputs.HasLine("Done installing git-hooks."),
		outputs.HasLine("Executing setup commands."),
		outputs.HasLine("Execute setup command: "+anotherCommand),
		outputs.HasLine(anotherCommand+" executed successfully"),
		outputs.HasLine("Execute setup command: "+someCommand),
		outputs.HasLine(someCommand+" executed successfully"),
		outputs.HasLine("Done executing setup commands."),
	)
}

func Test_IfSetupIsRunWithoutSkippingGitHooks_SetupAll_ShouldNotPrintSkipMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()

	output, _ := f.MakeMrtCommand().
		Setup().
		All().
		Execute()

	output.AssertHasNoLineContaining(t, "Skipping install-git-hooks step.")
}

func Test_IfSetupCommandExistsWithoutSkipping_SetupAll_ShouldNotPrintSkipMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()
	commandName := "some-command"
	f.SetupFixture.WriteSpyCommand(commandName)

	output, _ := f.MakeMrtCommand().
		Setup().
		All().
		Execute()

	output.AssertHasNoLineContaining(t, "Skipping setup command: "+commandName)
}

func Test_IfSetupIsRunWithSkipCloneRepositories_SetupAll_ShouldNotCloneRepositories(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()
	repositoryName := "1_TestRepository"
	f.TeamConfigWriter().Write(teamconfig.WithRepositories([]string{git.MakeCloneURL(repositoryName)}))

	f.MakeMrtCommand().
		Setup().
		All("--skip-clone-repositories").
		Execute()

	f.AssertFolderDoesNotExist("repositories/" + repositoryName)
}

func Test_IfSetupIsRunWithSkipCloneRepositories_SetupAll_ShouldPrintSkipMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()
	repositoryName := "1_TestRepository"
	f.TeamConfigWriter().Write(teamconfig.WithRepositories([]string{git.MakeCloneURL(repositoryName)}))

	output, _ := f.MakeMrtCommand().
		Setup().
		All("--skip-clone-repositories").
		Execute()

	output.AssertInOrder(t, outputs.HasLine("Skipping clone-repositories step."))
}

func Test_IfSetupIsRunWithSkipGitHooks_SetupAll_ShouldPrintSkipMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()

	output, _ := f.MakeMrtCommand().
		Setup().
		All("--skip-install-git-hooks").
		Execute()

	output.AssertHasLine(t, "Skipping install-git-hooks step.")
}

func Test_IfTwoSetupCommandsExistAndFirstIsSkipped_SetupAll_ShouldOnlyRunSecond(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	someCommand := "some-command"
	anotherCommand := "another-command"
	f.SetupFixture.WriteSpyCommand(someCommand)
	f.SetupFixture.WriteSpyCommand(anotherCommand)

	f.MakeMrtCommand().
		Setup().
		All("--skip-" + someCommand).
		Execute()

	f.SetupFixture.AssertSpyWasNotCalled(t, someCommand)
	f.SetupFixture.AssertSpyWasCalledWith(t, anotherCommand, f.SetupFixture.RepoDir)
}

func Test_IfTwoSetupCommandsExistAndSecondIsSkipped_SetupAll_ShouldOnlyRunFirst(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	someCommand := "some-command"
	anotherCommand := "another-command"
	f.SetupFixture.WriteSpyCommand(someCommand)
	f.SetupFixture.WriteSpyCommand(anotherCommand)

	f.MakeMrtCommand().
		Setup().
		All("--skip-" + anotherCommand).
		Execute()

	f.SetupFixture.AssertSpyWasCalledWith(t, someCommand, f.SetupFixture.RepoDir)
	f.SetupFixture.AssertSpyWasNotCalled(t, anotherCommand)
}

func Test_IfSetupCommandIsSkipped_SetupAll_ShouldPrintSkipMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	commandName := "some-command"
	f.SetupFixture.WriteSpyCommand(commandName)

	output, _ := f.MakeMrtCommand().
		Setup().
		All("--skip-" + commandName).
		Execute()

	output.AssertHasLine(t, "Skipping setup command: "+commandName)
}
