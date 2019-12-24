package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

// Exec runs a command with a pseudo-tty and returns an
// error if the program exits with a non-zero code.
func Exec(command string, args ...string) error {
	cmd := exec.Command(command, args...)

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
