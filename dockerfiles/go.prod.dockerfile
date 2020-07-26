# This is a generic Dockerfile used for running golang services in production. It's referenced in the deployment config file.

FROM golang:1.14-alpine

WORKDIR /app
COPY . .

RUN go mod download

ARG service_name
RUN CGO_ENABLED=0 GOOS=linux go install ./${service_name}

FROM alpine:latest

# In order for a build argument to be available in the CMD, we must make it an
# environment variable. This is because the CMD is only executed at runtime.
# The ARG command must be after FROM to be available at this point in the Dockerfile.
ARG service_name
ENV SERVICE ${service_name}

EXPOSE 80
WORKDIR /root/
COPY --from=0 /go/bin/${service_name} .
CMD ["./${SERVICE}"]
