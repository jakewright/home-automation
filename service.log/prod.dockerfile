FROM golang:latest

WORKDIR /go/src/home-automation
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go install ./service.log

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/bin/service.log .
CMD ["./service.log"]
