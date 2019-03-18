FROM golang:latest
RUN go get -u golang.org/x/lint/golint
RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /go/src/home-automation
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD CompileDaemon -build="go install ./service.config" -command="/go/bin/service.config"