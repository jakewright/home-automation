package watch

import (
	"context"
	"sync"
	"time"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/services/log/domain"
	"github.com/jakewright/home-automation/services/log/repository"

	"github.com/fsnotify/fsnotify"
)

// Watcher notifies subscribers of new events whenever the log file is written to
type Watcher struct {
	// LogDAO provides access to the log events
	LogRepository *repository.LogRepository

	subscribers map[chan<- *domain.Event]*repository.LogQuery
	mux         sync.Mutex        // Concurrent map access
	notify      chan struct{}     // Triggers reading new events from the log files
	ticker      *time.Ticker      // Used as a rate limiter
	watcher     *fsnotify.Watcher // Internal file watcher
}

// GetName returns the name "watcher"
func (w *Watcher) GetName() string {
	return "watcher"
}

// Start begins watching for log file changes and notifies subscribers accordingly
func (w *Watcher) Start() error {
	// Make sure the receiver struct has been initialised properly
	if w.LogRepository == nil {
		return oops.InternalService("LogRepository is not set")
	}
	if w.LogRepository.LogDirectory == "" {
		return oops.InternalService("Log directory is not set")
	}

	// Create an fsnotify watcher and attach to w so
	// that the Stop method can call Close() on it
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return oops.WithMessage(err, "failed to create file watcher")
	}
	defer func() { _ = watcher.Close() }()
	w.watcher = watcher

	// Start watching the log file directory so we
	// are notified when new log files are created
	if err = watcher.Add(w.LogRepository.LogDirectory); err != nil {
		return oops.WithMessage(err, "failed to watch log directory")
	}
	slog.Infof("Watching %s for changes", w.LogRepository.LogDirectory)

	// Create a notification channel with a buffer of 1 so
	// that we can always queue a new event while the current
	// one is in process. If the channel was unbuffered, we would
	// risk missing events if the notifier were not ready to
	// receive when the file write happened.
	w.notify = make(chan struct{}, 1)

	// Create a ticker to act as the rate limiter when notifying
	// subscribers. Without this then we risk thrashing the disk.
	w.ticker = time.NewTicker(time.Second * 2)

	go w.notifySubscribers()

	for {
		select {
		case fileEvent, ok := <-watcher.Events:
			if !ok {
				// If the channel is closed then just exit silently
				// because Stop() was probably called
				return nil
			}

			// We'll get a write event if any file inside the directory is written to.
			// If the file isn't actually a log file we'll waste some work
			// trying to read new events but it's safe to do.
			if fileEvent.Op&fsnotify.Write != fsnotify.Write {
				continue
			}

			// Write to the notify channel but do not block.
			// If the channel is not ready to receive then skip.
			select {
			case w.notify <- struct{}{}:
			default:
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				// If the channel is closed then just exit silently
				// because Stop() was probably called
				return nil
			}

			// It's unclear what state the watcher will be in if we receive
			// any errors so just return, which will trigger Close()
			return oops.WithMessage(err, "received error from watcher")
		}
	}
}

// Stop stops watching for log file changes
func (w *Watcher) Stop(_ context.Context) error {
	w.ticker.Stop()
	close(w.notify)

	if w.watcher != nil {
		return w.watcher.Close()
	}

	return nil
}

// Subscribe starts sending all events that match the query over the given channel. The query
// will be updated with the a new SinceUUID value whenever events are published to the channel.
func (w *Watcher) Subscribe(c chan<- *domain.Event, q *repository.LogQuery) error {
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

// notifySubscribers finds and sends new events to all subscribers
// whenever the notify channel is written to, rate limited by w.ticker.
func (w *Watcher) notifySubscribers() {
	// Read notify events until the channel is closed
	for range w.notify {
		// Block on the ticker to rate limit
		<-w.ticker.C

		w.findAndSendEvents()
	}
}

// findAndSendEvents will find all new events for each subscriber
// and send them over the subscribers' channels
func (w *Watcher) findAndSendEvents() {
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
			slog.Errorf("Failed to get events for subscriber: %v", err)
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
