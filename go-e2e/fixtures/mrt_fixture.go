package fixtures

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"

	"mrt-cli/go-e2e/assert"
	"mrt-cli/go-e2e/git"
	"mrt-cli/go-e2e/internal"
	mrtclient "mrt-cli/go-e2e/mrt"
	"mrt-cli/go-e2e/ssh"
	"mrt-cli/go-e2e/teamconfig"
)

type MrtFixture struct {
	t            *testing.T
	binaryPath   string
	agent        *ssh.Agent
	tempDir      string
	identityFile string
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

	return &MrtFixture{
		t:            t,
		binaryPath:   getBinaryPath(internal.GetRepoRoot(), t),
		agent:        agent,
		tempDir:      t.TempDir(),
		identityFile: "/dev/null",
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

func (f *MrtFixture) GitClone(repositoryName string, destination string) {
	git.MakeCommand(f.isolatedEnv()).
		Clone(git.MakeCloneURL(repositoryName), f.tempDir+"/"+destination).
		Execute()
}

func (f *MrtFixture) MakeMrtCommand() mrtclient.DirectedCommand {
	return mrtclient.
		MakeCommand(f.binaryPath, f.isolatedEnv()).
		RunInDirectory(f.tempDir)
}

func (f *MrtFixture) isolatedEnv() []string {
	sshConfigPath := internal.GetRepoRoot() + "/.ssh/config"
	sshCommand := "ssh -F " + sshConfigPath + " -o IdentityFile=" + f.identityFile + " -o IdentitiesOnly=yes"
	return append(f.agent.Env(), "GIT_SSH_COMMAND="+sshCommand)
}

func (f *MrtFixture) WriteTeamJSON(withOptions ...teamconfig.Option) {
	teamconfig.WriteTo(f.tempDir, withOptions...)
}

func (f *MrtFixture) AssertRepositoryExists(repositoryName string, inFolder string) {
	f.t.Helper()
	assert.DirectoryExists(f.t, f.tempDir+"/"+inFolder+"/"+repositoryName+"/.git")
}

func (f *MrtFixture) AssertFolderDoesNotExist(folder string) {
	f.t.Helper()
	assert.DirectoryDoesNotExist(f.t, f.tempDir+"/"+folder)
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
