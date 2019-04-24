package watch

import (
	"context"
	"home-automation/libraries/go/errors"
	"home-automation/libraries/go/slog"
	"sync"

	"github.com/fsnotify/fsnotify"
	"home-automation/service.log/domain"
	"home-automation/service.log/repository"
)

// Watcher notifies subscribers of new events whenever the log file is written to
type Watcher struct {
	// LogDAO provides access to the log events
	LogRepository *repository.LogRepository

	// Location is the path to the log file to watch
	Location string

	watcher     *fsnotify.Watcher
	subscribers map[chan<- *domain.Event]*repository.LogQuery
	mux         sync.Mutex
}

// GetName returns the name "watcher"
func (w *Watcher) GetName() string {
	return "watcher"
}

// Start begins watching for log file changes and notifies subscribers accordingly
func (w *Watcher) Start() error {
	metadata := map[string]string{
		"location": w.Location,
	}

	if w.LogRepository == nil {
		return errors.InternalService("LogRepository is not set", metadata)
	}

	if w.Location == "" {
		return errors.InternalService("File location is not set", metadata)
	}

	// Create a file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, metadata)
	}
	defer watcher.Close()
	w.watcher = watcher

	// Start watching the log file
	err = watcher.Add(w.Location)
	if err != nil {
		return errors.Wrap(err, metadata)
	}
	slog.Info("Watching log file for changes", metadata)

	for {
		select {
		case fileEvent, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if fileEvent.Op&fsnotify.Write != fsnotify.Write {
				continue
			}

			w.notifySubscribers()
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}

			return errors.Wrap(err, metadata)
		}
	}
}

// Stop stops watching for log file changes
func (w *Watcher) Stop(ctx context.Context) error {
	if w.watcher != nil {
		return w.watcher.Close()
	}

	return nil
}

// Subscribe starts sending all events that match the query over the given channel. The query
// will be updated with the a new SinceUUID value whenever events are published to the channel.
func (w *Watcher) Subscribe(c chan<- *domain.Event, q *repository.LogQuery) error {
	if q.SinceUUID == "" {
		return errors.InternalService("SinceUUID not set in subscriber query")
	}

	// Obtain a lock so we can write to the map
	w.mux.Lock()
	defer w.mux.Unlock()

	// Initialise the map if necessary
	if w.subscribers == nil {
		w.subscribers = make(map[chan<- *domain.Event]*repository.LogQuery)
	}

	// A channel is comparable so it's fine to use as a key
	w.subscribers[c] = q

	return nil
}

// Unsubscribe stops publishing events to the channel but does not close the channel
func (w *Watcher) Unsubscribe(c chan<- *domain.Event) {
	w.mux.Lock()
	defer w.mux.Unlock()
	delete(w.subscribers, c)
}

func (w *Watcher) notifySubscribers() {
	// Obtain a write lock before doing anything so that
	// we don't send duplicate events to the subscriber
	w.mux.Lock()
	defer w.mux.Unlock()

	for c, q := range w.subscribers {
		// Ensure that events are always published in order
		q.Reverse = false

		// Get all new events for this subscriber
		events, err := w.LogRepository.Find(q)
		if err != nil {
			slog.Error("Failed to get events for subscriber: %v", err)
			continue
		}

		// Send the events over the channel
		for _, event := range events {
			select {
			case c <- event: // Non-blocking write to the channel
			default: // Don't log otherwise we get a cycle of logs
			}
		}

		// Update the query for this subscriber
		if len(events) > 0 {
			// Events will always be in order so we can take the UUID of the last one
			q.SinceUUID = events[len(events)-1].UUID
		}
	}
}
