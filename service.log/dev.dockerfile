FROM golang:1.13
RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /app
COPY . .

RUN go get -v -t -d ./...

# Must use exec form so that CompileDaemon receives signals. The graceful-kill option then forwards them to the go binary.
CMD ["CompileDaemon", "-build=go install ./service.log", "-command=/go/bin/service.log", "-log-prefix=false", "-graceful-kill=true", "-graceful-timeout=10"]
