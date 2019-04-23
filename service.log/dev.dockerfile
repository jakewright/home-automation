FROM golang:latest
RUN go get github.com/githubnemo/CompileDaemon
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/home-automation
COPY . .

RUN dep ensure

# Must use exec form so that CompileDaemon receives signals. The graceful-kill option then forwards them to the go binary.
CMD ["CompileDaemon", "-build=go install ./service.log", "-command=/go/bin/service.log", "-log-prefix=false", "-graceful-kill=true"]
