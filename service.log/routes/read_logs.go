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
	Services string
	Severity int
	Since    string // The HTML datetime-local element formats time weirdly so we need to unmarshal to a string
	Until    string
	Reverse  bool
}

type readLogsResponse struct {
	FormattedEvents []*domain.FormattedEvent
	Services        string
	Severity        int
	Since           string
	Until           string
	Reverse         bool
}

const htmlTimeFormat = "2006-01-02T15:04"

func HandleReadLogs(w http.ResponseWriter, r *http.Request) {
	body := readLogsRequest{}
	if err := request.Decode(r, &body); err != nil {
		response.WriteJSON(w, err)
		return
	}

	// Default to logs from the last hour
	var err error
	since := time.Now().Add(-1 * time.Hour)
	until := time.Now()

	if body.Since != "" {
		since, err = time.Parse(htmlTimeFormat, body.Since)
		if err != nil {
			response.WriteJSON(w, err)
			return
		}
	}

	if body.Until != "" {
		until, err = time.Parse(htmlTimeFormat, body.Until)
		if err != nil {
			response.WriteJSON(w, err)
			return
		}
	}

	var services []string
	if body.Services != "" {
		services = strings.Split(strings.Replace(body.Services, " ", "", -1), ",")
	}

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

		formattedEvents = append(formattedEvents, event.Format())
	}

	// This is counter-intuitive but it is correct
	if !body.Reverse {
		reverse(formattedEvents)
	}

	rsp := readLogsResponse{
		FormattedEvents: formattedEvents,
		Services:        body.Services,
		Severity:        body.Severity,
		Since:           since.Format(htmlTimeFormat),
		Until:           until.Format(htmlTimeFormat),
		Reverse:         body.Reverse,
	}

	t, err := template.ParseFiles("service.log/templates/index.html")
	if err != nil {
		slog.Error("Failed to parse template: %v", err)
		response.WriteJSON(w, err)
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, rsp)
	if err != nil {
		slog.Error("Failed to execute template: %v", err)
		response.WriteJSON(w, err)
		return
	}

	response.Write(w, buf)

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

func reverse(a []*domain.FormattedEvent) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}
