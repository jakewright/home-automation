package repository

import (
	"bytes"
	"home-automation/libraries/go/errors"
	"home-automation/libraries/go/slog"
	"io/ioutil"
	"strings"
	"time"

	"home-automation/service.log/domain"
)

// LogRepository provides a query interface to the log file
type LogRepository struct {
	// Location is the path to the log file
	Location string
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
	lines, err := r.readLines()
	if err != nil {
		return nil, err
	}

	var events []*domain.Event

	// Iterate backwards so we process newer log lines first
	for i := len(lines) - 1; i >= 0; i-- {
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
			break
		}

		// Filter by UUID
		if q.SinceUUID != "" && event.UUID == q.SinceUUID {
			break
		}

		events = append(events, event)
	}

	// This is counter-intuitive but it is correct
	if !q.Reverse {
		reverse(events)
	}

	return events, nil
}

// readLines loads all lines from the log file into memory
func (r *LogRepository) readLines() ([][]byte, error) {
	data, err := ioutil.ReadFile(r.Location)
	if err != nil {
		return nil, errors.Wrap(err, nil)
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
