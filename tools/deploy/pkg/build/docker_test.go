package build

import "testing"

func Test_compareDockerfileArgs(t *testing.T) {
	tests := []struct {
		name              string
		dockerfileContent string
		givenArgs         map[string]string
		wantErr           bool
	}{
		{
			name:              "No args required",
			dockerfileContent: egDockerfileNoArgs,
			givenArgs:         nil,
			wantErr:           false,
		},
		{
			name:              "Missing both args",
			dockerfileContent: egDockerfileArgs,
			givenArgs:         nil,
			wantErr:           true,
		},
		{
			name:              "Missing one arg",
			dockerfileContent: egDockerfileArgs,
			givenArgs:         map[string]string{"work_dir": "/"},
			wantErr:           true,
		},
		{
			name:              "Got both args",
			dockerfileContent: egDockerfileArgs,
			givenArgs:         map[string]string{"work_dir": "/", "service_name": "foo"},
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := compareDockerfileArgs(tt.dockerfileContent, tt.givenArgs); (err != nil) != tt.wantErr {
				t.Errorf("compareDockerfileArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

const egDockerfileNoArgs = `FROM golang
COPY . .
`

const egDockerfileArgs = `FROM golang

ARG work_dir
WORKDIR ${work_dir}

COPY . .
ARG service_name
RUN go install ./${service_name}
`
