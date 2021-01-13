// shell is a package implementing a simple `sh` inspired shell with builtins
// It supports a very small subset of what proper shells like sh, ash, or bash
// have, but is still useful as a basic interactive console to run commands.
package shell

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type Handler func(ctx context.Context, args []string, sh *Shell) (int, error)

// ExitError from a handler causes the shell to exit
// Used by the `exit` builtin.
type ExitError struct{}
func (e ExitError) Error() string {
	return "exit"
}

// a Shell is the active instance of a shell, created with shell.New
type Shell struct {
	in io.Reader
	out io.Writer

	// Environment variables
	env map[string]string

	// Built-in handlers
	handlers map[string]Handler

	// Current working directory
	cwd string
}

// New creates a shell with the given input and outputs.
func New(input io.Reader, output io.Writer) *Shell {
	return &Shell{
		in:  input,
		out: output,

		env: map[string]string{
			"PS1": "$ ",
		},

		handlers: map[string]Handler{
			"exit": Exit,
			"cd": CD,
			"echo": Echo,
		},
	}
}

// Run an interactive shell: Prompt for commands, execute them, loop.
// Blocks until completion or the context is cancelled
func (sh *Shell) Interact(ctx context.Context) error {
	for {
		// Write shell prompt:
		if _, err := sh.out.Write([]byte(sh.env["PS1"])); err != nil {
			return err
		}

		scan := bufio.NewScanner(sh.in)
		// Read input command:
		if !scan.Scan() {
			return scan.Err()
		}

		// Parse command in a totally janky way
		scanner := bufio.NewScanner(bytes.NewReader(scan.Bytes()))
		scanner.Split(bufio.ScanWords)
		if !scanner.Scan() {
			continue
		}
		cmd := scanner.Text()

		args := []string{cmd}
		for scanner.Scan() {
			args = append(args, scanner.Text())
		}

		var ret int
		var err error
		if handler, ok := sh.handlers[cmd]; ok {
			ret, err = handler(ctx, args, sh)
		} else {
			ret, err = sh.Exec(ctx, args)
		}

		if err != nil {
			// Built-ins can cause us to exit:
			if errors.Is(err, ExitError{}) {
				return nil
			} else {
				sh.out.Write([]byte(fmt.Sprintf("error running command: %v", err)))
			}
		}

		// Store return value as $?
		sh.env["?"] = fmt.Sprintf("%d", ret)
	}
}

func (sh *Shell) Exec(ctx context.Context, args []string) (int, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = sh.in
	cmd.Stdout = sh.out
	cmd.Stderr = sh.out
	cmd.Dir = sh.cwd
	err := cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError); if ok {
			return exitErr.ExitCode(), nil
		}
		// todo return codes
		_, _ = sh.out.Write([]byte(fmt.Sprintf("error %v\n", err)))
		return 1, nil
	}
	return 0, nil
}
