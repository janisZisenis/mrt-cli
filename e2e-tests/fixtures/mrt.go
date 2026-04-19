package fixtures

import (
	"context"
	"mrt-cli/e2e-tests/assert"
	"mrt-cli/e2e-tests/git"
	"mrt-cli/e2e-tests/internal"
	"mrt-cli/e2e-tests/ssh"
	"mrt-cli/e2e-tests/teamconfig"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	mrtclient "mrt-cli/e2e-tests/mrt"
)

const setupCommandDir = "setup"

type MrtFixture struct {
	t            *testing.T
	binaryPath   string
	agent        *ssh.Agent
	TeamDir      string
	identityFile string
	RunFixture   *RunCommandFixture
	SetupFixture *CommandFixture
}

func MakeMrtFixture(t *testing.T) *MrtFixture {
	t.Helper()
	t.Parallel()

	agent, err := ssh.StartAgent()
	if err != nil {
		t.Fatalf("Could not start SSH agent.")
	}

	t.Cleanup(func() {
		stopErr := agent.Stop()
		if stopErr != nil {
			t.Fatalf("%v", stopErr)
		}
	})

	// EvalSymlinks is needed on macOS where t.TempDir() returns a path under
	// /var/folders which is a symlink to /private/var/folders. Without this,
	// assertions comparing TeamDir against paths returned by the OS would fail.
	teamDir, _ := filepath.EvalSymlinks(t.TempDir())

	return &MrtFixture{
		t:            t,
		binaryPath:   getBinaryPath(internal.GetRepoRoot(), t),
		agent:        agent,
		TeamDir:      teamDir,
		identityFile: "/dev/null",
		RunFixture:   NewRunCommandFixture(teamDir),
		SetupFixture: NewCommandFixture(teamDir, setupCommandDir),
	}
}

func (f *MrtFixture) Authenticate() *MrtFixture {
	f.t.Helper()
	privateKeyPath := internal.GetRepoRoot() + "/.ssh/private-key"
	f.identityFile = privateKeyPath

	if err := f.agent.AddKey(privateKeyPath); err != nil {
		f.t.Fatalf("%v", err)
	}

	return f
}

func (f *MrtFixture) AbsolutePath(relativePath string) string {
	return f.TeamDir + "/" + relativePath
}

func (f *MrtFixture) MakeGitCommand() git.BaseCommand {
	return git.MakeCommand(f.isolatedEnv())
}

func (f *MrtFixture) MakeMrtCommand() mrtclient.BaseCommand {
	return mrtclient.MakeCommand(f.binaryPath, f.isolatedEnv())
}

func (f *MrtFixture) MakeMrtCommandInTeamDir() mrtclient.DirectedCommand {
	return mrtclient.
		MakeCommand(f.binaryPath, f.isolatedEnv()).
		RunInDirectory(f.TeamDir)
}

func (f *MrtFixture) isolatedEnv() []string {
	sshConfigPath := internal.GetRepoRoot() + "/.ssh/config"
	sshCommand := "ssh -F " + sshConfigPath + " -o IdentityFile=" + f.identityFile + " -o IdentitiesOnly=yes"
	binaryDir := filepath.Dir(f.binaryPath)
	return append(f.agent.Env(),
		"GIT_SSH_COMMAND="+sshCommand,
		"PATH="+os.Getenv("PATH")+":"+binaryDir,
	)
}

func (f *MrtFixture) TeamConfigWriter() *teamconfig.Writer {
	return teamconfig.NewWriter(f.TeamDir)
}

func (f *MrtFixture) AssertRepositoryExists(repositoryName string, inFolder string) {
	f.t.Helper()
	assert.DirectoryExists(f.t, f.TeamDir+"/"+inFolder+"/"+repositoryName+"/.git")
}

func (f *MrtFixture) AssertFolderDoesNotExist(folder string) {
	f.t.Helper()
	assert.DirectoryDoesNotExist(f.t, f.TeamDir+"/"+folder)
}

func getBinaryPath(repositoryDir string, t *testing.T) string {
	cmd := exec.CommandContext(
		context.Background(),
		"mrt",
		"--team-dir",
		repositoryDir,
		"run",
		"binary-location",
	)
	binaryPathBytes, err := cmd.Output()

	output := string(binaryPathBytes)

	if err != nil {
		t.Fatalf("failed to run mrt command: %s", output)
	}

	binaryPath := stringTrimNewline(output)

	if binaryPath == "" {
		t.Fatalf("command returned empty directory")
	}

	_, statErr := os.Stat(binaryPath)
	if os.IsNotExist(statErr) {
		t.Fatalf("binary not found at: %s", binaryPath)
	}

	return binaryPath
}

func stringTrimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return strings.TrimSpace(s)
}
