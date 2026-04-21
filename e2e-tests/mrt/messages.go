package mrt

import "fmt"

const (
	MsgFailedToLoadTeamConfiguration = "Failed to load team configuration"
	MsgFailedToCloneRepository       = "Failed to clone repository, skipping it."

	MsgInvalidRepositoriesPath = "repositoriesPath must be a relative path within the team repository"

	MsgInstallingGitHooksToTeamRepository     = "Installing git-hooks to team repository"
	MsgDoneInstallingGitHooksToTeamRepository = "Done installing git-hooks to team repository."
)

func MsgActionNotAllowedOnBranch(action, branch string) string {
	return fmt.Sprintf("Action %q not allowed on branch %q", action, branch)
}

func MsgInstallingGitHooksToRepositoriesLocatedIn(dir string) string {
	return fmt.Sprintf("Installing git-hooks to repositories located in %q", dir)
}

func MsgCloning(url string) string {
	return "Cloning " + url
}

func MsgSuccessfullyCloned(url string) string {
	return "Successfully cloned " + url
}
