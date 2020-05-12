package exe

import (
	"bytes"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

// RemoteCmd is a command that is executed over SSH
type RemoteCmd struct {
	Cmd       string
	PseudoTTY bool
}

// RequestPseudoTTY will force the command to request a pty on the remote server.
func (c *RemoteCmd) RequestPseudoTTY() *RemoteCmd {
	c.PseudoTTY = true
	return c
}

// RemoteCommand creates a new remote command. You can specify the whole command
// as a single string, or as separate arguments that get joined by spaces.
func RemoteCommand(name string, args ...string) *RemoteCmd {
	elems := append([]string{name}, args...)
	return &RemoteCmd{
		Cmd: strings.Join(elems, " "),
	}
}

// Run executes the command using the given SSH client
func (c *RemoteCmd) Run(client *ssh.Client) Result {
	errParams := map[string]string{
		"cmd": c.Cmd,
	}

	session, err := client.NewSession()
	if err != nil {
		return Result{
			Err: oops.WithMessage(err, "failed to start session", errParams),
		}
	}
	defer func() { _ = session.Close() }()

	if c.PseudoTTY {
		modes := ssh.TerminalModes{
			ssh.ECHO: 0, // disable echoing
		}
		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			return Result{
				Err: oops.WithMessage(err, "failed to request Pty", errParams),
			}
		}
	}

	result := Result{}

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(c.Cmd)

	result.Stdout = strings.TrimSpace(stdout.String())
	result.Stderr = strings.TrimSpace(stderr.String())

	if err != nil {
		result.Err = oops.WithMessage(err, result.Stderr, errParams)
	}

	return result
}

// SSHClient returns an SSH client for the given username@host.
// It assumes ssh-agent will be available for key authentication.
func SSHClient(username, host string) (*ssh.Client, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to open SSH_AUTH_SOCK")
	}

	agentClient := agent.NewClient(conn)

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", clientConfig)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to dial %s", host)
	}

	return client, nil
}
