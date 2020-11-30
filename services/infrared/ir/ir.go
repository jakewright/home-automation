package ir

import (
	"context"
	"time"

	"github.com/jakewright/home-automation/libraries/go/distsync"
	lircproxydef "github.com/jakewright/home-automation/services/lirc-proxy/def"
)

type Instruction func(context.Context, lircproxydef.LircProxyService) error

func Wait(ms int) Instruction {
	return func(context.Context, lircproxydef.LircProxyService) error {
		time.Sleep(time.Millisecond * time.Duration(ms))
		return nil
	}
}

func Key(device, key string) Instruction {
	return func(ctx context.Context, lirc lircproxydef.LircProxyService) error {
		return send(ctx, device, key)
	}
}

type IRSend struct {
	LIRC lircproxydef.LircProxyService
}

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

func send(
	ctx context.Context,
	lirc lircproxydef.LircProxyService,
	device, key string,
) error {
	if _, err := lirc.SendOnce(ctx, &lircproxydef.SendOnceRequest{
		Device: device,
		Key:    key,
	}).Wait(); err != nil {
		return
	}
}
