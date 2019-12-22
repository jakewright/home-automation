package systemd

import (
	"bytes"
	"log"
	"text/template"

	"github.com/jakewright/home-automation/tools/deploy/config"
)

func deployGo(service *config.Service) error {
	// Checkout the code locally
	// Build the binary
	// scp the binary to the target
	// SSH to the target
	// Stop, create, restart service

	dir := service.Target.Directory
	cmd := "" // @todo

	data := struct {
		PreStop, StopService, PostStop, CreateAndStartService string
	}{
		StopService:           stopService(service.Name),
		CreateAndStartService: createAndStartService(service.Name, config.LangGo, dir, cmd),
	}

	txt := `
ssh -t -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null {{ .Username }}@{{ .Host }} << EOF
	# Abort if anything fails
	set -e

	{{ StopService }}

	{{ CreateAndStartService }}
EOF
`
	tmpl, err := template.New("ssh").Parse(txt)
	if err != nil {
		log.Fatal(err)
	}

	b := bytes.Buffer{}
	if err := tmpl.Execute(&b, data); err != nil {
		log.Fatal(err)
	}

	return b.String()
}
