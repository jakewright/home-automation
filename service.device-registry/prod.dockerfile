FROM golang:latest
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/jakewright/home-automation
COPY . .

RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go install ./service.device-registry

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/bin/service.device-registry .
CMD ["./service.log"]
