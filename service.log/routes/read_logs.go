package routes

import (
	"bytes"
	"home-automation/libraries/go/response"
	"html/template"
	"io/ioutil"
	"net/http"

	"home-automation/service.log/domain"
)

func HandleReadLogs(w http.ResponseWriter, r *http.Request) {
	// Read all log messages
	b, err := ioutil.ReadFile("/var/log/messages")
	if err != nil {
		response.WriteJSON(w, err)
		return
	}

	// Take the latest 10 lines
	bLines := bytes.Split(b, []byte("\n"))
	bLines = bLines[len(bLines)-10:]

	lines := make([]*domain.Line, len(bLines))

	// Format each line
	for i, b := range bLines {
		lines[i] = domain.NewLineFromBytes(b)
	}

	log := domain.Log{
		Lines: lines,
	}

	t, err := template.ParseFiles("service.log/templates/index.html")
	if err != nil {
		response.WriteJSON(w, err)
		return
	}

	err = t.Execute(w, log)
	if err != nil {
		response.WriteJSON(w, err)
		return
	}
}
