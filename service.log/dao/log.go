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

type LogRepository struct {
	// Location is the path to the log file
	Location string

	Events chan *domain.Event

	// lineCount is the number of lines in the log file so the watcher knows what has changed
	lineCount int
}

func NewLogRepository(location string) *LogRepository {
	return &LogRepository{
		Location: location,
	}
}

func (r *LogRepository) Find(services []string, severity slog.Severity, since, until time.Time) ([]*domain.Event, error) {
	lines, err := r.readLines()
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

func (r *LogRepository) Watch() error {
	metadata := map[string]string{
		"location": r.Location,
	}

	// To only emit events for new log lines, we need to store
	// how many lines currently exist in the log file.
	lines, err := r.readLines()
	if err != nil {
		return err
	}
	r.lineCount = len(lines)

	// Create a file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, nil)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := r.sendNewLines(); err != nil {
						slog.Error("Failed to send new lines: %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("Error watching log file: %v", err, metadata)
			}
		}
	}()

	err = watcher.Add(r.Location)
	if err != nil {
		slog.Panic("Failed to add log to watcher: %v", err)
	}

	<-done
	return nil
}

func (r *LogRepository) readLines() ([][]byte, error) {
	data, err := ioutil.ReadFile(r.Location)
	if err != nil {
		return nil, errors.Wrap(err, nil)
	}

	return bytes.Split(data, []byte("\n")), nil
}

func (r *LogRepository) sendNewLines() error {
	lines, err := r.readLines()
	if err != nil {
		return err
	}

	lineCount := len(lines)
	var start, end int

	switch {
	// File is empty
	case lineCount == 0:
		r.lineCount = 0
		return nil

	// Nothing has changed
	case lineCount == r.lineCount:
		return nil

	// Probably started a new log file
	case lineCount < r.lineCount:
		start = 0
		end = lineCount - 1

	// New log lines
	case lineCount > r.lineCount:
		start = int(math.Max(float64(r.lineCount-1), 0))
		end = lineCount
	}

	for i := start; i < end; i++ {
		// Skip empty lines
		if len(lines[i]) == 0 {
			continue
		}

		r.Events <- domain.NewEventFromBytes(i, lines[i])
	}

	r.lineCount = lineCount
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
