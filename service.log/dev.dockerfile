FROM golang:latest
RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /go/src/home-automation
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD CompileDaemon -build="go install ./service.log" -command="/go/bin/service.log" -log-prefix=false
