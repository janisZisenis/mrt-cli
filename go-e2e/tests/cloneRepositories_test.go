package tests

import (
	"bytes"
	"fmt"
	"mrt-cli/go-e2e/assertions"
	"mrt-cli/go-e2e/fixtures"

	"mrt-cli/go-e2e/utils"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCloneRepositoriesToCustomPath(t *testing.T) {
	t.Parallel()
	repositoryDir := "/Users/jazi/Documents/Development/mrt/repositories/mrt-cli"

	agent := fixtures.AuthenticatedFixture(t, repositoryDir+"/.ssh/private-key")
	_ = setupPath(repositoryDir)
	_ = agent.ShowKeys()
	mrtCmd := exec.Command("mrt", "--team-dir", repositoryDir, "run", "binary-location", "--", "--exe-name")
	exePathBytes, err := mrtCmd.Output()
	if err != nil {
		t.Fatalf("failed to get executable path: %v", err)
	}
	exePath := stringTrimNewline(string(exePathBytes))
	tempDir := t.TempDir()

	teamConfig := tempDir + "/team.json"
	repositoryName := "1_TestRepository"
	_ = utils.WriteJSONFile(teamConfig, map[string]interface{}{
		"repositories": []string{"git@github-testing:janisZisenisTesting/" + repositoryName + ".git"},
	})

	cmd := exec.Command(exePath, "--team-dir", tempDir, "setup", "clone-repositories")
	cmd.Env = append(os.Environ(), agent.Env()...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	assertions.TestDirectoryExists(t, tempDir+"/repositories/"+repositoryName+"/.git")
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
