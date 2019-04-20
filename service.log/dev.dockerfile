FROM golang:latest
RUN go get github.com/githubnemo/CompileDaemon
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/home-automation
COPY . .

RUN dep ensure

CMD CompileDaemon -build="go install ./service.log" -command="/go/bin/service.log" -log-prefix=false
