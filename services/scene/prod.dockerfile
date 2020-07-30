FROM golang:1.13

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install ./service.scene

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/bin/service.scene .
CMD ["./service.scene"]
