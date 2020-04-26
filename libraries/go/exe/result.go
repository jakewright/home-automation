package exe

// Result represents the outcome of running a command.
type Result struct {
	Err    error
	Stdout string
	Stderr string
}
