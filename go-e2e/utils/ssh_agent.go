package utils

import (
	"context"
	"errors"
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

type SSHAddKeyError struct {
	KeyPath string
	Err     error
	Output  string
}

func (e *SSHAddKeyError) Error() string {
	return fmt.Sprintf("failed to add key %s: %v\noutput: %s", e.KeyPath, e.Err, e.Output)
}

func (a *Agent) AddKey(pathToKey string) error {
	cmd := exec.CommandContext(context.Background(), "ssh-add", pathToKey)
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &SSHAddKeyError{
			KeyPath: pathToKey,
			Err:     err,
			Output:  string(out),
		}
	}
	return nil
}

type SSHRemoveKeyError struct {
	KeyPath string
	Err     error
	Output  string
}

func (e *SSHRemoveKeyError) Error() string {
	return fmt.Sprintf("failed to remove key %s: %v\noutput: %s", e.KeyPath, e.Err, e.Output)
}

func (a *Agent) RemoveKey(pathToKey string) error {
	cmd := exec.CommandContext(context.Background(), "ssh-add", "-d", pathToKey)
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &SSHRemoveKeyError{
			KeyPath: pathToKey,
			Err:     err,
			Output:  string(out),
		}
	}
	return nil
}

type SSHStartError struct {
	Err error
}

func (e *SSHStartError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("failed to start ssh-agent: %v", e.Err)
	}
	return "failed to start ssh-agent"
}

func StartSSHAgent() (*Agent, error) {
	output, err := exec.CommandContext(context.Background(), "ssh-agent", "-s").Output()
	if err != nil {
		return nil, &SSHStartError{Err: err}
	}

	sock := parseVariable(output, "SSH_AUTH_SOCK")
	pid := parseVariable(output, "SSH_AGENT_PID")

	if sock == nil || pid == nil {
		return nil, &SSHStartError{
			Err: errors.New("failed to parse ssh-agent output"),
		}
	}

	return &Agent{Sock: *sock, PID: *pid}, nil
}

const sshAgentOutputParts = 2

func parseVariable(out []byte, varName string) *string {
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, varName+"=") {
			parts := strings.SplitN(line, ";", sshAgentOutputParts)
			value := strings.TrimPrefix(parts[0], varName+"=")
			return &value
		}
	}
	return nil
}

type SSHStopError struct {
	PID string
	Err error
}

func (e *SSHStopError) Error() string {
	return fmt.Sprintf("failed to stop ssh-agent PID=%s: %v", e.PID, e.Err)
}

func (a *Agent) Stop() error {
	cmd := exec.CommandContext(context.Background(), "ssh-agent", "-k")
	cmd.Env = append(cmd.Env, a.Env()...)
	if err := cmd.Run(); err != nil {
		return &SSHStopError{
			PID: a.PID,
			Err: err,
		}
	}
	return nil
}

type SSHShowKeysError struct {
	Err    error
	Output string
}

func (e *SSHShowKeysError) Error() string {
	return fmt.Sprintf("failed to show keys: %v\noutput: %s", e.Err, e.Output)
}

func (a *Agent) ShowKeys() error {
	cmd := exec.CommandContext(context.Background(), "ssh-add", "-l")
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &SSHShowKeysError{
			Err:    err,
			Output: string(out),
		}
	}

	return nil
}
