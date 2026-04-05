package fixtures

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"mrt-cli/go-e2e/assert"
	"mrt-cli/go-e2e/git"
	"mrt-cli/go-e2e/internal"
	mrtclient "mrt-cli/go-e2e/mrt"
	"mrt-cli/go-e2e/ssh"
	"mrt-cli/go-e2e/teamconfig"
)

const setupCommandDir = "setup"

type MrtFixture struct {
	t            *testing.T
	binaryPath   string
	agent        *ssh.Agent
	teamDir      string
	identityFile string
	RunFixture   *RunCommandFixture
	SetupFixture *CommandFixture
}

func MakeMrtFixture(t *testing.T) *MrtFixture {
	t.Helper()

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

	teamDir := t.TempDir()

	return &MrtFixture{
		t:            t,
		binaryPath:   getBinaryPath(internal.GetRepoRoot(), t),
		agent:        agent,
		teamDir:      teamDir,
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

func (f *MrtFixture) Parallel() *MrtFixture {
	f.t.Helper()
	f.t.Parallel()

	return f
}

func (f *MrtFixture) AbsolutePath(relativePath string) string {
	return f.teamDir + "/" + relativePath
}

func (f *MrtFixture) TeamDir() string {
	return f.teamDir
}

func (f *MrtFixture) MakeGitCommand() git.BaseCommand {
	return git.MakeCommand(f.isolatedEnv())
}

func (f *MrtFixture) MakeMrtCommand() mrtclient.DirectedCommand {
	return mrtclient.
		MakeCommand(f.binaryPath, f.isolatedEnv()).
		RunInDirectory(f.teamDir)
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
	return teamconfig.NewWriter(f.teamDir)
}

func (f *MrtFixture) AssertRepositoryExists(repositoryName string, inFolder string) {
	f.t.Helper()
	assert.DirectoryExists(f.t, f.teamDir+"/"+inFolder+"/"+repositoryName+"/.git")
}

func (f *MrtFixture) AssertFolderDoesNotExist(folder string) {
	f.t.Helper()
	assert.DirectoryDoesNotExist(f.t, f.teamDir+"/"+folder)
}

func getBinaryPath(repositoryDir string, t *testing.T) string {
	cmd := exec.CommandContext(context.Background(), "mrt", "--team-dir", repositoryDir, "run", "binary-location")
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
