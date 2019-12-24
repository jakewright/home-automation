package systemd

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/jakewright/home-automation/tools/deploy/config"
)

func createSSHScript(serviceName, language, workingDirectory, cmd string) string {
	return ""
}

func stopService(serviceName string) string {
	// The "or true" will suppress the non-zero exit code if the service does not
	// exist. Note that this will still print an error to the console though.
	return fmt.Sprintf("sudo systemctl stop %s || true\n", unitFileName(serviceName))
}

func createAndStartService(serviceName, language, workingDirectory, cmd string) string {
	unit := unit(serviceName, language, workingDirectory, cmd)

	data := struct {
		UnitFileName, Unit string
	}{
		UnitFileName: unitFileName(serviceName),
		Unit:         unit,
	}

	txt := `
sudo cat >/lib/systemd/system/{{ .UnitFileName }} <<EOL
{{ .Unit }}
EOL

sudo chmod 644 /lib/systemd/system/{{ .UnitFileName}}

sudo systemctl daemon-reload
sudo systemctl enable {{ .UnitFileName }}
sudo systemctl start {{ .UnitFileName }}
`

	tmpl, err := template.New("systemd").Parse(txt)
	if err != nil {
		log.Fatal(err)
	}

	b := bytes.Buffer{}
	if err := tmpl.Execute(&b, data); err != nil {
		log.Fatal(err)
	}

	return b.String()
}

func unit(serviceName, language, workingDirectory, cmd string) string {
	switch language {
	case config.LangGo:
		return goUnit(serviceName, workingDirectory, cmd)
	}

	log.Fatalf("no systemd unit definition for language '%s'\n", language)
	return ""
}

func goUnit(serviceName, workingDirectory, cmd string) string {
	data := struct {
		ServiceName, SyslogIdentifier, WorkingDirectory, ExecStart string
	}{
		ServiceName:      serviceName,
		SyslogIdentifier: syslogIdentifier(serviceName),
		WorkingDirectory: workingDirectory,
		ExecStart:        cmd,
	}

	txt := `[Unit]
Description={{ .ServiceName }}

[Service]
SyslogIdentifier={{ .SyslogIdentifier }}
WorkingDirectory={{ .WorkingDirectory }}
Type=idle
ExecStart={{ .ExecStart }}

[Install]
WantedBy=multi-user.target`

	tmpl, err := template.New("unit").Parse(txt)
	if err != nil {
		log.Fatal(err)
	}

	b := bytes.Buffer{}
	if err := tmpl.Execute(&b, data); err != nil {
		log.Fatal(err)
	}

	return b.String()
}

func unitFileName(serviceName string) string {
	serviceNameDashes := strings.Replace(serviceName, ".", "-", -1)
	return serviceNameDashes + ".service"
}

func syslogIdentifier(serviceName string) string {
	serviceNameDashes := strings.Replace(serviceName, ".", "-", -1)

	// The syslog identifier must start with ha- so that it's matched by
	// the rsyslog rule that forwards the logs to a central syslog server.
	return "ha-" + serviceNameDashes
}
