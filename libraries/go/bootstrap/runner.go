package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jakewright/home-automation/libraries/go/slog"
)

// Process is a long-running task that provides service functionality
type Process interface {
	// GetName returns a friendly name for the process for use in logs
	GetName() string

	// Start kicks off the task and only returns when the task has finished.
	// The task will be stopped if the context is cancelled.
	Start(ctx context.Context) error
}

type runner struct {
	processes []Process
	deferred  []func() error
	mu        sync.Mutex
	running   bool
}

func (r *runner) addProcess(p Process) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		panic("cannot add process after Run() called")
	}

	r.processes = append(r.processes, p)
}

func (r *runner) addDeferred(d func() error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		panic("cannot add deferred func after Run() called")
	}

	r.deferred = append(r.deferred, d)
}

// Run takes a number of processes and concurrently runs them all. It will stop
// if all processes terminate or if a signal (SIGINT or SIGTERM) is received.
func (r *runner) Run() {
	r.mu.Lock()
	r.running = true
	r.mu.Unlock()

	// os.Exit should be the last thing to happen
	var code int
	defer os.Exit(code)

	// Close all of the resources after processes have shut down
	for _, deferred := range r.deferred {
		defer func(d func() error) {
			if err := d(); err != nil {
				code = 1
			}
		}(deferred)
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	// Start all of the processes in goroutines
	for _, process := range r.processes {
		process := process

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := process.Start(ctx); err != nil {
				slog.Errorf("Process %s stopped with error: %v", process.GetName(), err)
				code = 1
			} else {
				slog.Debugf("Process %s stopped", process.GetName())
			}
		}()
	}

	// Close the done channel when all processes return
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for all processes to return or for a signal
	select {
	case <-done: // All processes stopped
	case s := <-sig:
		slog.Infof("Received %v signal", s)
	}

	// Cancelling the context will signal to all processes
	// that they should shutdown. If all processes have
	// already stopped, this is a no-op.
	cancel()

	// Wait for processes to terminate
	wg.Wait()
	slog.Infof("All processes stopped")
}
