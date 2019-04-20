package routes

import (
	"bytes"
	"home-automation/libraries/go/request"
	"home-automation/libraries/go/response"
	"home-automation/libraries/go/slog"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"home-automation/service.log/domain"
)

type readLogsRequest struct {
	Since    time.Time
	Until    time.Time
	Severity int
	Services string
}

func HandleReadLogs(w http.ResponseWriter, r *http.Request) {
	body := readLogsRequest{
		// Default to logs from the last hour
		Since: time.Now().Add(-1 * time.Hour),
		Until: time.Now(),
	}

	if err := request.Decode(r, &body); err != nil {
		slog.Error("Failed to decode request: %v", err)
		response.WriteJSON(w, err)
		return
	}

	var services []string
	if body.Services != "" {
		services = strings.Split(body.Services, ",")
	}

	slog.Info("Since: %v", body.Since.String())
	slog.Info("Services: %v", services)

	// Read all log messages
	data, err := ioutil.ReadFile("/var/log/messages")
	if err != nil {
		slog.Error("Failed to read logs: %v", err)
		response.WriteJSON(w, err)
		return
	}

	lines := bytes.Split(data, []byte("\n"))
	var formattedEvents []*domain.FormattedEvent

	// Start at len - 2 because the last line is always an empty line
	for i := len(lines) - 2; i >= 0; i-- {
		event := domain.NewEventFromBytes(i, lines[i])

		// Filter by severity
		if int(event.Severity) < body.Severity {
			continue
		}

		// Filter by services
		if len(services) > 0 {
			if !contains(services, event.Service) {
				continue
			}
		}

		// Filter by time
		if event.Timestamp.After(body.Until) {
			continue
		}
		if event.Timestamp.Before(body.Since) {
			break
		}

		formattedEvents = append(formattedEvents, event.Format())
	}

	reverse(formattedEvents)

	log := domain.Log{
		FormattedEvents: formattedEvents,
	}

	t, err := template.ParseFiles("service.log/templates/index.html")
	if err != nil {
		slog.Error("Failed to parse template: %v", err)
		response.WriteJSON(w, err)
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, log)
	if err != nil {
		slog.Error("Failed to execute template: %v", err)
		response.WriteJSON(w, err)
		return
	}

	response.Write(w, buf)

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func reverse(a []*domain.FormattedEvent) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}
