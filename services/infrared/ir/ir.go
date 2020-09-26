package ir

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/jakewright/home-automation/libraries/go/distsync"
	"github.com/jakewright/home-automation/libraries/go/oops"
)

type Instruction func(context.Context) error

func Wait(ms int) Instruction {
	return func(context.Context) error {
		time.Sleep(time.Millisecond * time.Duration(ms))
		return nil
	}
}

func Key(device, key string) Instruction {
	return func(ctx context.Context) error {
		return send(ctx, device, key)
	}
}

type IRSend struct{}

func (s *IRSend) Execute(ctx context.Context, ins []Instruction) error {
	// dsync will take care of time outs
	// and context cancellations for us
	lock, err := distsync.Lock(ctx, "ir")
	if err != nil {
		return err
	}
	defer lock.Unlock()

	for _, instruction := range ins {
		if err := instruction(ctx); err != nil {
			return err
		}
	}

	return nil
}

func send(ctx context.Context, device, key string) error {
	cmd := exec.CommandContext(ctx, "irsend", "SEND_ONCE", device, key)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return oops.WithMessage(err, "failed to run `irsend SEND_ONCE %s %s`", device, key)
	}

	// Assume anything written to stderr is a bad thing
	if stderr.Len() > 0 {
		return oops.InternalService("irsend wrote to stderr: %s", stderr.String())
	}

	return nil
}
