package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	tick  = "\xE2\x9C\x94"
	green = "\033[32m"
	reset = "\033[0m"
)

// BuildDirectory is injected at compile time
var BuildDirectory string

const usage = `Home Automation Runner
USAGE
	Start, stop or restart a set of services
	  run [stop|restart] service.name group

	Build a service or set of services
	  build service.name group

	Apply all schema files for a set of services (default all)
	  schema apply [service.name] group

	Seed mock data for a set of services (default all)
	  schema seed [service.name] group

GROUPS
	core
	  Core platform services

	log
	  Services required for logging
`

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
	case "build":
		build(os.Args[2:])
	case "schema":
		schema(os.Args[2:])
	default:
		start(os.Args[1:])
	}
}

func help() {
	log.Print(usage)
	os.Exit(0)
}

func start(args []string) {
	var build bool
	if len(args) > 0 && args[0] == "--build" {
		build = true
		args = args[1:]
	}

	services := getServices(args)
	composeArgs := []string{"up", "-d", "--renew-anon-volumes", "--remove-orphans"}
	if build {
		composeArgs = append(composeArgs, "--build")
	}
	composeArgs = append(composeArgs, services...)

	if len(services) == 0 {
		log.Printf("Starting all services...\n")
	} else {
		log.Printf("Starting %s...\n", strings.Join(services, ", "))
	}

	runPTY("docker-compose", composeArgs)
}

func stop(args []string) {
	services := getServices(args)
	composeArgs := append([]string{"stop"}, getServices(args)...)

	if len(services) == 0 {
		log.Printf("Stopping all services...\n")
	} else {
		log.Printf("Stopping %s...\n", strings.Join(services, ", "))
	}

	runPTY("docker-compose", composeArgs)
}

func restart(args []string) {
	stop(args)
	start(args)
}

func build(args []string) {
	services := getServices(args)
	if len(services) == 0 {
		log.Printf("Nothing to build")
		return
	}

	log.Printf("Building %s...\n", strings.Join(services, ", "))

	composeArgs := append([]string{"build", "--pull"}, services...)
	runPTY("docker-compose", composeArgs)
}

func schema(args []string) {
	if len(args) < 1 {
		log.Fatal(usage)
	}

	services := getServices(args[1:])
	if len(services) == 0 {
		services = getAllServiceNames()
	}

	running := false

	// If the MySQL container exists
	containerID := run("docker-compose", "ps", "-q", "mysql")
	if containerID != "" {
		// See if it's running
		result := run("docker", "inspect", "-f", "{{.State.Running}}", containerID)
		b, err := strconv.ParseBool(result)
		if err != nil {
			log.Fatal(err)
		}

		running = b
	}

	if !running {
		start([]string{"mysql"})
		containerID = run("docker-compose", "ps", "-q", "mysql")
		time.Sleep(time.Second * 5)
	}

	switch args[0] {
	case "apply":
		applySchema(services, containerID)
	case "seed":
		seedData(services, containerID)
	default:
		log.Fatal(usage)
	}
}

func applySchema(services []string, containerID string) {
	log.Printf("Applying schema...\n")

	for _, name := range services {
		schema := getServiceSchema(name)
		if schema == "" {
			continue
		}

		fmt.Printf(name)

		args := []string{"exec", "-i", containerID, "sh", "-c", "exec mysql -uroot -psecret"}
		runWithInput(schema, "docker", args...)

		fmt.Printf(" %s%s%s\n", green, tick, reset)
	}
}

func seedData(services []string, containerID string) {
	log.Printf("Seeding data...\n")

	for _, name := range services {
		sql := getServiceMockSQL(name)
		if sql == "" {
			continue
		}

		fmt.Printf(name)

		args := []string{"exec", "-i", containerID, "sh", "-c", "exec mysql -uroot -psecret"}
		runWithInput(sql, "docker", args...)

		fmt.Printf(" %s%s%s\n", green, tick, reset)
	}
}

func getAllServiceNames() []string {
	var services []string

	ls, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range ls {
		if !info.IsDir() {
			continue
		}

		if !strings.HasPrefix(info.Name(), "service.") {
			continue
		}

		services = append(services, info.Name())
	}

	return services
}

func getServiceSchema(service string) string {
	filename := "./" + service + "/schema/schema.sql"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return ""
	} else if err != nil {
		log.Fatal(err)
	}

	a := "USE home_automation;\n\n"

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return a + string(b)
}

func getServiceMockSQL(service string) string {
	filename := "./" + service + "/schema/mock_data.sql"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return ""
	} else if err != nil {
		log.Fatal(err)
	}

	a := "USE home_automation;\n\n"

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return a + string(b)
}

func getServices(args []string) []string {
	var services []string
	for _, s := range args {
		services = append(services, expandService(s)...)
	}
	return services
}

func expandService(s string) []string {
	coreServices := []string{"service.api-gateway", "service.config", "service.device-registry", "redis"}
	logServices := []string{"filebeat", "logstash", "service.log"}

	switch s {
	case "core":
		return coreServices
	case "log":
		return logServices
	}

	return []string{s}
}

func run(command string, args ...string) string {
	return runWithInput("", command, args...)
}

func runWithInput(stdin, command string, args ...string) string {
	cmd := exec.Command(command, args...)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer func() { _ = pipe.Close() }()
		_, _ = io.WriteString(pipe, stdin)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\n%s\n", out)
		log.Fatal(err)
	}

	return string(bytes.TrimSpace(out))
}

func runPTY(command string, args []string) {
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
