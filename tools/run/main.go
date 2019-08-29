package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

// BuildDirectory is injected at compile time
var BuildDirectory string

const usage = "Usage: run [stop|restart] core log service.name"

func main() {
	if cwd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else if cwd != BuildDirectory {
		log.Fatalf("Must be run from home-automation root: %s\n", BuildDirectory)
	}

	if len(os.Args) < 2 {
		log.Fatal(usage)
	}

	switch os.Args[1] {
	case "help", "--help":
		help()
	case "stop":
		stop(os.Args[2:])
	case "restart":
		restart(os.Args[2:])
	default:
		start(os.Args[1:])
	}
}

func help() {
	log.Println(usage)
	os.Exit(0)
}

func start(args []string) {
	services := getServices(args)
	composeArgs := append([]string{"up", "-d", "--renew-anon-volumes", "--remove-orphans"}, services...)

	if len(services) == 0 {
		log.Printf("Starting all services...\n")
	} else {
		log.Printf("Starting %s...\n", strings.Join(services, ", "))
	}

	run("docker-compose", composeArgs)
}

func stop(args []string) {
	services := getServices(args)
	composeArgs := append([]string{"stop"}, getServices(args)...)

	if len(services) == 0 {
		log.Printf("Stopping all services...\n")
	} else {
		log.Printf("Stopping %s...\n", strings.Join(services, ", "))
	}

	run("docker-compose", composeArgs)
}

func restart(args[]string) {
	stop(args)
	start(args)
}

func getServices(args []string) []string {
	var services []string
	for _, s := range args {
		services = append(services, expandService(s)...)
	}
	return services
}

func expandService(s string) []string {
	coreServices := []string{"service.api-gateway", "service.config", "service.registry.device", "redis"}
	logServices := []string{"filebeat", "logstash", "service.log"}

	switch s {
	case "core":
		return coreServices
	case "log":
		return logServices
	}

	return []string{s}
}

func run(command string, args []string) {
	cmd := exec.Command(command, args...)

	// Start the command with a pty
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
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
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort

	// Copy stdin to the terminal and the terminal to stdout
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)
}