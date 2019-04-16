package routes

import (
	"bytes"
	"home-automation/libraries/go/response"
	"home-automation/libraries/go/slog"
	"html/template"
	"io/ioutil"
	"net/http"

	"home-automation/service.log/domain"
)

type Controller struct {
}

func HandleReadLogs(w http.ResponseWriter, r *http.Request) {
	// Read all log messages
	b, err := ioutil.ReadFile("/var/log/messages")
	if err != nil {
		slog.Error("Failed to read logs: %v", err)
		response.WriteJSON(w, err)
		return
	}

	// Take the latest 10 lines
	bLines := bytes.Split(b, []byte("\n"))
	bLines = bLines[len(bLines)-11 : len(bLines)-1]

	lines := make([]*domain.FormattedLine, len(bLines))

	// Format each line
	for i, b := range bLines {
		lines[i] = domain.NewLineFromBytes(b).FormatLine()
	}

	log := domain.Log{
		FormattedLines: lines,
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
