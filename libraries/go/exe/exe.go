package exe

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

// Cmd represents an external command being prepared or run.
type Cmd struct {
	*exec.Cmd
	PseudoTTY bool
	Input     string
}

// SetPseudoTTY is a fluent setter for the pty option.
func (c *Cmd) SetPseudoTTY() *Cmd {
	c.PseudoTTY = true
	return c
}

// SetInput is a fluent setter for the input option.
func (c *Cmd) SetInput(input string) *Cmd {
	c.Input = input
	return c
}

// Dir is a fluent setter for the directory option.
func (c *Cmd) Dir(dir string) *Cmd {
	c.Cmd.Dir = dir
	return c
}

// Env is a fluent setter for the env option.
// Elements should be of the form NAME=value.
func (c *Cmd) Env(env []string) *Cmd {
	c.Cmd.Env = env
	return c
}

// Command returns a new Cmd.
func Command(name string, args ...string) *Cmd {
	return &Cmd{
		Cmd: exec.Command(name, args...),
	}
}

// Run runs the command and returns a Result
func (c *Cmd) Run() *Result {
	errParams := map[string]string{
		"cmd": c.String(),
	}

	// TODO: support the other options when PseudoTTY is set
	if c.PseudoTTY {
		if err := runPTY(c.Cmd); err != nil {
			return &Result{
				Err: oops.WithMetadata(err, errParams),
			}
		}

		return &Result{}
	}

	if c.Input != "" {
		pipe, err := c.Cmd.StdinPipe()
		if err != nil {
			return &Result{Err: err}
		}

		go func() {
			defer func() { _ = pipe.Close() }()
			_, _ = io.WriteString(pipe, c.Input)
		}()
	}

	var stdout, stderr bytes.Buffer
	c.Cmd.Stdout = &stdout
	c.Cmd.Stderr = &stderr

	result := &Result{}

	err := c.Cmd.Run()

	result.Stdout = strings.TrimSpace(stdout.String())
	result.Stderr = strings.TrimSpace(stderr.String())

	if err != nil {
		result.Err = oops.WithMessage(err, result.Stderr, errParams)
	}

	return result
}

func runPTY(cmd *exec.Cmd) error {
	// Start the command with a pty
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	// Make sure to close the terminal at the end
	defer func() { _ = ptmx.Close() }() // Best effort

	// Handle terminal size
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Fatalf("Failed to resize pty: %v", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize

	// Put terminal into raw mode
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort

	// Copy stdin to the terminal and the terminal to stdout
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	// Wait for the command to finish to avoid nil pointer exception below
	if err := cmd.Wait(); err != nil {
		return err
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("non-zero exit code: %d", cmd.ProcessState.ExitCode())
	}

	return nil
}
