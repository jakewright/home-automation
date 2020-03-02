FROM golang:1.13
RUN go get github.com/jakewright/compile-daemon

WORKDIR /app
COPY . .

RUN go get -v -t -d ./...

# Must use exec form so that compile-daemon receives signals. The graceful-kill option then forwards them to the go binary.
CMD ["compile-daemon", "-build=go install ./service.log", "-command=/go/bin/service.log", "-directories=service.log,libraries/go", "-log-prefix=false", "-graceful-kill=true", "-graceful-timeout=10"]
