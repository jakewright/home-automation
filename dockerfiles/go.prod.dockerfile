# This is a generic Dockerfile used for running golang services in production.
# It's referenced in the deployment config file.

FROM golang:1.15-alpine

WORKDIR /app
COPY . .

RUN go mod download

ARG service_name
ARG revision
RUN CGO_ENABLED=0 GOOS=linux go install \
    -ldflags "-X github.com/jakewright/home-automation/libraries/go/bootstrap.Revision=${revision}" \
    ./services/${service_name}

FROM alpine:latest

# In order for a build argument to be available in the CMD, we must make it an
# environment variable. This is because the CMD is only executed at runtime.
# The ARG command must be after FROM to be available at this point in the Dockerfile.
ARG service_name
ENV SERVICE ${service_name}

EXPOSE 80
WORKDIR /root/
COPY --from=0 /go/bin/${service_name} .

# Copy assets for this service into /assets in the image. The LICENCE file is
# included because Docker requires at least one file in the COPY command, and
# it is assumed that LICENCE will always exist. Any file could be used.
COPY LICENCE ./private/assets/${service_name}/prod/* /assets/

# Use the shell form of CMD so that the environment variable gets executed
CMD ./${SERVICE}
