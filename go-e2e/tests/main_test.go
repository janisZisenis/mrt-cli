package tests

import (
	"mrt-cli/go-e2e/utils"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binaryName string

func TestMain(m *testing.M) {
	repositoryDir := utils.GetRepoRootDir()
	binaryDir := getBinaryDir(repositoryDir)
	binaryName = getBinaryName(repositoryDir)

	panicIfBinaryCanNotBeFound(binaryDir, binaryName)

	addToPath(binaryDir)

	os.Exit(m.Run())
}

func panicIfBinaryCanNotBeFound(binaryDir string, binaryName string) {
	binaryPath := binaryDir + "/" + binaryName
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		panic("binary not found at: " + binaryPath)
	}
}

func addToPath(binaryDir string) {
	currentPath := os.Getenv("PATH")
	newPath := currentPath + string(os.PathListSeparator) + binaryDir

	if err := os.Setenv("PATH", newPath); err != nil {
		panic("failed to update add binaryDir to PATH: " + binaryDir)
	}
}

func getBinaryDir(repositoryDir string) string {
	cmd := exec.Command("mrt", "--team-dir", repositoryDir, "run", "binary-location", "--", "--dir")
	binaryDirBytes, err := cmd.Output()

	output := string(binaryDirBytes)

	if err != nil {
		panic("failed to run mrt command: " + output)
	}

	binaryDir := strings.TrimSpace(output)

	if binaryDir == "" {
		panic("command returned empty directory")
	}

	return binaryDir
}

func getBinaryName(repositoryDir string) string {
	cmd := exec.Command("mrt", "--team-dir", repositoryDir, "run", "binary-location", "--", "--exe-name")
	exePathBytes, err := cmd.Output()

	output := string(exePathBytes)

	if err != nil {
		panic("failed to get executable path: " + output)
	}

	return stringTrimNewline(output)
}

func stringTrimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}
