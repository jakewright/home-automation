package routes

import (
	"bytes"
	"home-automation/libraries/go/request"
	"home-automation/libraries/go/response"
	"home-automation/libraries/go/slog"
	"html/template"
	"net/http"
	"strings"
	"time"

	"home-automation/service.log/dao"

	"home-automation/service.log/domain"
)

type Controller struct {
	Repository *dao.LogRepository
}

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

func (c *Controller) HandleReadLogs(w http.ResponseWriter, r *http.Request) {
	body := readLogsRequest{}
	if err := request.Decode(r, &body); err != nil {
		response.WriteJSON(w, err)
		return
	}

	var services []string
	if body.Services != "" {
		services = strings.Split(strings.Replace(body.Services, " ", "", -1), ",")
	}

	severity := slog.Severity(body.Severity)

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

	metadata := map[string]string{
		"services": body.Services,
		"severity": severity.String(),
		"since":    since.Format(time.RFC3339),
		"until":    since.Format(time.RFC3339),
	}

	events, err := c.Repository.Find(services, severity, since, until)
	if err != nil {
		slog.Error("Failed to find events: %v", err, metadata)
		response.WriteJSON(w, err)
		return
	}

	formattedEvents := make([]*domain.FormattedEvent, len(events))
	for i, event := range events {
		formattedEvents[i] = event.Format()
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

func reverse(a []*domain.FormattedEvent) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}
