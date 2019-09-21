package config

import (
	"context"
	"time"

	"github.com/jakewright/home-automation/libraries/go/api"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

const (
	maxRetries     = 5
	backoff        = time.Second * 2
	reloadInterval = time.Second * 30
)

type Loader struct {
	ServiceName string

	config   *Config
	ticker   *time.Ticker
	done     chan struct{}
	reloaded time.Time
}

func (l *Loader) GetName() string {
	return "config"
}

func (l *Loader) Load() error {
	var content map[string]interface{}
	var err error
	for i := 0; i < maxRetries; i++ {
		if _, err = api.Get("service.config/read/"+l.ServiceName, &content); err == nil {
			break
		}
		slog.Error("Failed to load config [attempt %d of %d]: %v", i+1, maxRetries, err)
		time.Sleep(backoff)
	}

	if err != nil {
		return err
	}

	if l.config == nil {
		l.config = New(content)
		DefaultProvider = l.config
		return nil
	}

	l.config.Replace(content)
	return nil
}

func (l *Loader) Start() error {
	l.done = make(chan struct{})
	l.ticker = time.NewTicker(reloadInterval)

	for {
		select {
		case <-l.done:
			return nil
		case <-l.ticker.C:
			if err := l.Load(); err == nil {
				l.reloaded = time.Now()
			}
		}
	}
}

func (l *Loader) Stop(ctx context.Context) error {
	l.ticker.Stop()
	l.done <- struct{}{}
	return nil
}
