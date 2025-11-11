package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

type Agent struct {
	PID  string
	Sock string
}

func (a *Agent) Env() []string {
	return []string{
		"SSH_AGENT_PID=" + a.PID,
		"SSH_AUTH_SOCK=" + a.Sock,
	}
}

type SshAddKeyError struct {
	KeyPath string
	Err     error
	Output  string
}

func (e *SshAddKeyError) Error() string {
	return fmt.Sprintf("failed to add key %s: %v\noutput: %s", e.KeyPath, e.Err, e.Output)
}

func (a *Agent) AddKey(pathToKey string) error {
	cmd := exec.Command("ssh-add", pathToKey)
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &SshAddKeyError{
			KeyPath: pathToKey,
			Err:     err,
			Output:  string(out),
		}
	}
	return nil
}

type SshRemoveKeyError struct {
	KeyPath string
	Err     error
	Output  string
}

func (e *SshRemoveKeyError) Error() string {
	return fmt.Sprintf("failed to remove key %s: %v\noutput: %s", e.KeyPath, e.Err, e.Output)
}

func (a *Agent) RemoveKey(pathToKey string) error {
	cmd := exec.Command("ssh-add", "-d", pathToKey)
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &SshRemoveKeyError{
			KeyPath: pathToKey,
			Err:     err,
			Output:  string(out),
		}
	}
	return nil
}

type SshStartError struct {
	Err error
}

func (e *SshStartError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("failed to start ssh-agent: %v", e.Err)
	}
	return "failed to start ssh-agent"
}

func StartSSHAgent() (*Agent, error) {
	output, err := exec.Command("ssh-agent", "-s").Output()
	if err != nil {
		return nil, &SshStartError{Err: err}
	}

	sock := parseVariable(output, "SSH_AUTH_SOCK")
	pid := parseVariable(output, "SSH_AGENT_PID")

	if sock == nil || pid == nil {
		return nil, &SshStartError{
			Err: fmt.Errorf("failed to parse ssh-agent output"),
		}
	}

	return &Agent{Sock: *sock, PID: *pid}, nil
}

func parseVariable(out []byte, varName string) *string {
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, varName+"=") {
			parts := strings.SplitN(line, ";", 2)
			value := strings.TrimPrefix(parts[0], varName+"=")
			return &value
		}
	}
	return nil
}

type SshStopError struct {
	PID string
	Err error
}

func (e *SshStopError) Error() string {
	return fmt.Sprintf("failed to stop ssh-agent PID=%s: %v", e.PID, e.Err)
}

func (a *Agent) Stop() error {
	cmd := exec.Command("ssh-agent", "-k")
	cmd.Env = append(cmd.Env, a.Env()...)
	if err := cmd.Run(); err != nil {
		return &SshStopError{
			PID: a.PID,
			Err: err,
		}
	}
	return nil
}

type SshShowKeysError struct {
	Err    error
	Output string
}

func (e *SshShowKeysError) Error() string {
	return fmt.Sprintf("failed to show keys: %v\noutput: %s", e.Err, e.Output)
}

func (a *Agent) ShowKeys() error {
	cmd := exec.Command("ssh-add", "-l")
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &SshShowKeysError{
			Err:    err,
			Output: string(out),
		}
	}

	fmt.Println(string(out))
	return nil
}
