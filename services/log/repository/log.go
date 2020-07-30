package repository

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/services/log/domain"
)

// LogRepository provides a query interface to the log file
type LogRepository struct {
	// LogDirectory is the path to the directory containing daily log files
	LogDirectory string
}

// LogQuery is a set of conditions to apply when finding events
type LogQuery struct {
	// Services is a slice of service name patterns to filter by.
	// If the slice is empty, then events from all services will
	// be returned. Patterns may end with a wildcard "*" character.
	Services []string

	// Severity is the minimum severity that events need to have.
	// Set this to slog.Severity(0) to return all events.
	Severity slog.Severity

	// SinceTime is the earliest inclusive time that events should
	// be from. Set to the zero value to return all events.
	SinceTime time.Time

	// UntilTime is the latest inclusive time that events should
	// be from. Set to the zero value to return all events.
	UntilTime time.Time

	// SinceUUID is a UUID of an event. If not an empty string, only
	// events that happened _after_ this event will be returned. The
	// event with the given UUID itself will not be returned.
	SinceUUID string

	// Reverse will change the order of the returned results. If false,
	// events will be returned in chronological order, i.e. oldest first.
	Reverse bool
}

// Find returns all events that match the given query
func (r *LogRepository) Find(q *LogQuery) ([]*domain.Event, error) {
	events, err := r.findEvents(q)
	if err != nil {
		return nil, err
	}

	// This is counter-intuitive but it is correct
	if !q.Reverse {
		reverse(events)
	}

	return events, nil
}

func (r *LogRepository) findEvents(q *LogQuery) ([]*domain.Event, error) {
	var events []*domain.Event
	date := time.Now().UTC()

	for {
		filename := filepath.Join(r.LogDirectory, fmt.Sprintf("messages-%s", date.Format("2006-01-02")))

		lines, err := readLines(filename)
		if err != nil {
			// We expect to eventually find a file that does not exist so
			// don't return an error, just return the events found so far.
			if os.IsNotExist(err) {
				return events, nil
			}

			// Any other error is unexpected
			return nil, err
		}

		// Iterate backwards so we process newer log lines first
		for i := len(lines) - 1; i >= 0; i-- {
			// Skip empty lines
			if len(lines[i]) == 0 {
				continue
			}

			event := domain.NewEventFromBytes(lines[i])

			// Filter by severity
			if event.Severity < q.Severity {
				continue
			}

			// Filter by service
			if len(q.Services) > 0 && !containsService(q.Services, event.Service) {
				continue
			}

			// Filter by time
			if !q.UntilTime.IsZero() && event.Timestamp.After(q.UntilTime) {
				continue
			}
			if !q.SinceTime.IsZero() && event.Timestamp.Before(q.SinceTime) {
				return events, nil
			}

			// Filter by UUID
			if q.SinceUUID != "" && event.UUID == q.SinceUUID {
				return events, nil
			}

			events = append(events, event)
		}

		// Subtract a day from the date
		date = date.AddDate(0, 0, -1)
	}
}

// readLines loads all lines from the log file into memory
func readLines(filename string) ([][]byte, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to read log file")
	}

	return bytes.Split(data, []byte("\n")), nil
}

// containsService returns whether any of the patterns match the service name.
// Patterns may end with a wildcard character "*".
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

// reverse performs an in-place reversal of the given slice
func reverse(a []*domain.Event) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}
