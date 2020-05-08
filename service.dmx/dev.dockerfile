# Custom Dockerfile so an ola_set_dmx binary can be placed in the PATH

FROM golang:1.14-alpine

# Alpine doesn't have git but go get needs it
RUN apk add --no-cache git

# Use a fork of compile-daemon that supports watching multiple directories
RUN go get github.com/jakewright/compile-daemon

EXPOSE 80

COPY ./service.dmx/ola_set_dmx /bin/ola_set_dmx

WORKDIR /app
COPY . .

RUN go get -v -t -d ./...

# Must use exec form so that compile-daemon receives signals. The graceful-kill option then forwards them to the go binary.
# The -directories option doesn't work with the directories the other way around. It might be because of the dot in the service name.
CMD ["sh", "-c", "compile-daemon -build=\"go install ./service.dmx\" -command=/go/bin/service.dmx -directories=libraries/go,service.dmx -log-prefix=false -log-prefix=false -graceful-kill=true -graceful-timeout=10"]
