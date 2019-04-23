package dao

import (
	"bytes"
	"home-automation/libraries/go/errors"
	"home-automation/libraries/go/slog"
	"io/ioutil"
	"math"
	"strings"
	"time"

	"home-automation/service.log/domain"

	"github.com/fsnotify/fsnotify"
)

type LogDAO struct {
	// Location is the path to the log file
	Location string

	// Events is a channel over which all new log events are sent
	Events chan *domain.Event

	// lineCount is the number of lines in the log file so the watcher knows what has changed
	lineCount int
}

func NewLogRepository(location string) *LogDAO {
	return &LogDAO{
		Location: location,
		Events:   make(chan *domain.Event, 100),
	}
}

func (d *LogDAO) Find(services []string, severity slog.Severity, since, until time.Time) ([]*domain.Event, error) {
	lines, err := d.readLines()
	if err != nil {
		return nil, err
	}

	var events []*domain.Event

	// Iterate backwards so that newest log lines are at the front of the slice
	for i := len(lines) - 2; i >= 0; i-- {
		event := domain.NewEventFromBytes(i, lines[i])

		// Filter by severity
		if event.Severity < severity {
			continue
		}

		// Filter by service
		if len(services) > 0 {
			if !containsService(services, event.Service) {
				continue
			}
		}

		// Filter by time
		if event.Timestamp.After(until) {
			continue
		}
		if event.Timestamp.Before(since) {
			break
		}

		events = append(events, event)
	}

	return events, nil
}

func (d *LogDAO) Watch() {
	metadata := map[string]string{
		"location": d.Location,
	}

	// To only emit events for new log lines, we need to store
	// how many lines currently exist in the log file.
	lines, err := d.readLines()
	if err != nil {
		slog.Panic("Failed to read log lines: %v", err, metadata)
	}
	d.lineCount = len(lines)

	// Create a file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Panic("Failed to create file watcher: %v", err, metadata)
	}
	defer watcher.Close()

	err = watcher.Add(d.Location)
	if err != nil {
		slog.Panic("Failed to add log file to watcher: %v", err, metadata)
	}

	slog.Debug("Watching log file for changes", metadata)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				slog.Debug("Watcher events channel closed", metadata)
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := d.readNewEvents(); err != nil {
					slog.Error("Failed to read new events: %v", err, metadata)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				slog.Debug("Watcher errors channel closed", metadata)
				return
			}
			slog.Panic("Received error from file watcher: %v", err, metadata)
		}
	}
}

func (d *LogDAO) readLines() ([][]byte, error) {
	data, err := ioutil.ReadFile(d.Location)
	if err != nil {
		return nil, errors.Wrap(err, nil)
	}

	return bytes.Split(data, []byte("\n")), nil
}

// readNewEvents reads new events in the log file and sends each one over the Events channel
func (d *LogDAO) readNewEvents() error {
	lines, err := d.readLines()
	if err != nil {
		return err
	}

	lineCount := len(lines)
	var start, end int

	switch {
	// File is empty
	case lineCount == 0:
		d.lineCount = 0
		return nil

	// Nothing has changed
	case lineCount == d.lineCount:
		return nil

	// Probably started a new log file
	case lineCount < d.lineCount:
		start = 0
		end = lineCount - 1

	// New log lines
	case lineCount > d.lineCount:
		start = int(math.Max(float64(d.lineCount-1), 0))
		end = lineCount
	}

	for i := start; i < end; i++ {
		// Skip empty lines
		if len(lines[i]) == 0 {
			continue
		}

		d.Events <- domain.NewEventFromBytes(i, lines[i])
	}

	d.lineCount = lineCount
	return nil
}

func containsService(patterns []string, service string) bool {
	for _, p := range patterns {
		if p == service {
			return true
		}

		if p[len(p)-1:] == "*" && strings.HasPrefix(service, p[:len(p)-1]) {
			return true
		}
	}
	return false
}
