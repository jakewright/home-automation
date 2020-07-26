# This is a generic Dockerfile used for running golang services in production. It's referenced in the deployment config file.

FROM golang:1.14-alpine

WORKDIR /app
COPY . .

RUN go mod download

ARG service_name
RUN CGO_ENABLED=0 GOOS=linux go install ./${service_name}

FROM alpine:latest

EXPOSE 80
WORKDIR /root/

ARG service_name
COPY --from=0 /go/bin/${service_name} .

CMD ["./${service_name}"]
