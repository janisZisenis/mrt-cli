package ssh

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

type AddKeyError struct {
	KeyPath string
	Err     error
	Output  string
}

func (e *AddKeyError) Error() string {
	return fmt.Sprintf("failed to add key %s: %v\noutput: %s", e.KeyPath, e.Err, e.Output)
}

func (a *Agent) AddKey(pathToKey string) error {
	cmd := exec.CommandContext(context.Background(), "ssh-add", pathToKey)
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &AddKeyError{
			KeyPath: pathToKey,
			Err:     err,
			Output:  string(out),
		}
	}
	return nil
}

type RemoveKeyError struct {
	KeyPath string
	Err     error
	Output  string
}

func (e *RemoveKeyError) Error() string {
	return fmt.Sprintf("failed to remove key %s: %v\noutput: %s", e.KeyPath, e.Err, e.Output)
}

func (a *Agent) RemoveKey(pathToKey string) error {
	cmd := exec.CommandContext(context.Background(), "ssh-add", "-d", pathToKey)
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &RemoveKeyError{
			KeyPath: pathToKey,
			Err:     err,
			Output:  string(out),
		}
	}
	return nil
}

type StartError struct {
	Err error
}

func (e *StartError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("failed to start ssh-agent: %v", e.Err)
	}
	return "failed to start ssh-agent"
}

func StartAgent() (*Agent, error) {
	output, err := exec.CommandContext(context.Background(), "ssh-agent", "-s").Output()
	if err != nil {
		return nil, &StartError{Err: err}
	}

	sock := parseVariable(output, "SSH_AUTH_SOCK")
	pid := parseVariable(output, "SSH_AGENT_PID")

	if sock == nil || pid == nil {
		return nil, &StartError{
			Err: errors.New("failed to parse ssh-agent output"),
		}
	}

	return &Agent{Sock: *sock, PID: *pid}, nil
}

func parseVariable(out []byte, varName string) *string {
	agentOutputParts := 2
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, varName+"=") {
			parts := strings.SplitN(line, ";", agentOutputParts)
			value := strings.TrimPrefix(parts[0], varName+"=")
			return &value
		}
	}
	return nil
}

type StopError struct {
	PID string
	Err error
}

func (e *StopError) Error() string {
	return fmt.Sprintf("failed to stop ssh-agent PID=%s: %v", e.PID, e.Err)
}

func (a *Agent) Stop() error {
	cmd := exec.CommandContext(context.Background(), "ssh-agent", "-k")
	cmd.Env = append(cmd.Env, a.Env()...)
	if err := cmd.Run(); err != nil {
		return &StopError{
			PID: a.PID,
			Err: err,
		}
	}
	return nil
}

type ShowKeysError struct {
	Err    error
	Output string
}

func (e *ShowKeysError) Error() string {
	return fmt.Sprintf("failed to show keys: %v\noutput: %s", e.Err, e.Output)
}

func (a *Agent) ShowKeys() error {
	cmd := exec.CommandContext(context.Background(), "ssh-add", "-l")
	cmd.Env = append(cmd.Env, a.Env()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &ShowKeysError{
			Err:    err,
			Output: string(out),
		}
	}

	return nil
}
