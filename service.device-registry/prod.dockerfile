FROM golang:1.13

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install ./service.device-registry

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/bin/service.device-registry .
COPY ./private/devices/prod.json /data/config.json
CMD ["./service.device-registry"]
