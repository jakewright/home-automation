FROM golang:1.13

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install ./service.log

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/bin/service.log .
COPY ./service.log/templates /templates
CMD ["./service.log"]
