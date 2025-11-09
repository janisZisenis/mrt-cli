package tests

import (
	"bytes"
	"fmt"

	"mrt-cli/go-e2e/utils"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCloneRepositoriesToCustomPath(t *testing.T) {
	t.Parallel()

	//println()
	//val := os.Getenv("GIT_SSH_COMMAND")
	//if val == "" {
	//	t.Fatal("GIT_SSH_COMMAND not set")
	//}
	//t.Logf("GIT_SSH_COMMAND = %s", val)

	repositoryDir := "/Users/jazi/Documents/Development/mrt/repositories/mrt-cli"
	agent := utils.AuthenticatedFixture(t, repositoryDir+"/.ssh/private-key")
	_ = agent.ShowKeys()
	tempDir := t.TempDir()
	_ = setupPath(repositoryDir)
	teamConfig := tempDir + "/team.json"
	repositoryName := "1_TestRepository"
	_ = utils.WriteJSONFile(teamConfig, map[string]interface{}{
		"repositories": []string{"git@github-testing:janisZisenisTesting/" + repositoryName + ".git"},
	})

	mrtCmd := exec.Command("mrt", "--team-dir", repositoryDir, "run", "binary-location", "--", "--exe-name")
	exePathBytes, err := mrtCmd.Output()
	if err != nil {
		fmt.Printf("failed to get executable path: %v\n", err)
		os.Exit(1)
	}
	exePath := stringTrimNewline(string(exePathBytes))

	cmd := exec.Command(exePath, "--team-dir", tempDir, "setup", "clone-repositories")
	cmd.Env = append(os.Environ(), agent.Env()...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("command failed: %v\n", err)
		os.Exit(1)
	}

	utils.TestDirectoryExists(t, tempDir+"/repositories/"+repositoryName+"/.git")
}

func stringTrimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}

func setupPath(repositoryDir string) error {
	cmd := exec.Command("mrt", "--team-dir", repositoryDir, "run", "binary-location", "--", "--dir")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run mrt command: %w", err)
	}

	dir := strings.TrimSpace(out.String())

	if dir == "" {
		return fmt.Errorf("command returned empty directory")
	}

	currentPath := os.Getenv("PATH")
	newPath := currentPath + string(os.PathListSeparator) + dir

	if err := os.Setenv("PATH", newPath); err != nil {
		return fmt.Errorf("failed to update PATH: %w", err)
	}

	return nil
}
