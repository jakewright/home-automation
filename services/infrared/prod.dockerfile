FROM golang:1.13

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install ./service.infrared

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/bin/service.infrared .
COPY ./private/devices/prod.json /data/config.json
CMD ["./service.infrared"]
