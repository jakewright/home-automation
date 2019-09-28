FROM golang:1.13

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install ./service.config

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/bin/service.config .
COPY ./private/config/prod.yaml /data/config.yaml
CMD ["./service.config"]
