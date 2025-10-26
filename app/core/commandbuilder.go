package core

import (
	"context"
	"io"
	"os/exec"
	"time"
)

type CommandBuilder struct {
	command string
	args    []string
	stdout  io.Writer
	stderr  io.Writer
	stdin   io.Reader
	timeout time.Duration
}

const defaultCommandTimeout = 5

func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		stdout:  nil,
		stderr:  nil,
		stdin:   nil,
		timeout: defaultCommandTimeout * time.Second,
	}
}

func (b *CommandBuilder) WithCommand(command string) *CommandBuilder {
	b.command = command
	return b
}

func (b *CommandBuilder) WithArgs(args ...string) *CommandBuilder {
	b.args = append(b.args, args...)
	return b
}

func (b *CommandBuilder) WithStdout(stdout io.Writer) *CommandBuilder {
	b.stdout = stdout
	return b
}

func (b *CommandBuilder) WithStderr(stderr io.Writer) *CommandBuilder {
	b.stderr = stderr
	return b
}

func (b *CommandBuilder) WithStdin(stdin io.Reader) *CommandBuilder {
	b.stdin = stdin
	return b
}

func (b *CommandBuilder) Build() (*exec.Cmd, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)

	// #nosec G204 - We actually want the user to execute scripts and commands with their own arguments.
	// We don't know yet what these commands will be and what arguments they will pass to them.
	cmd := exec.CommandContext(ctx, b.command, b.args...)
	cmd.Stdout = b.stdout
	cmd.Stderr = b.stderr
	cmd.Stdin = b.stdin

	return cmd, ctx, cancel
}

func (b *CommandBuilder) Run() error {
	cmd, ctx, cancel := b.Build()
	defer cancel()

	err := cmd.Run()

	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	return err
}

func (b *CommandBuilder) Start() (context.CancelFunc, func() error, error) {
	cmd, ctx, cancel := b.Build()

	wait := func() error {
		err := cmd.Wait()

		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		return err
	}

	err := cmd.Start()
	return cancel, wait, err
}
