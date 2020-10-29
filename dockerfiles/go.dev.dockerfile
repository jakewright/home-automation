# This is a generic Dockerfile used for running golang services locally.
# It's referenced in the project's Docker Compose file.

FROM golang:1.15-alpine

# Alpine doesn't have git but go get needs it
RUN apk add --no-cache git

# Use a fork of compile-daemon that supports watching multiple directories
RUN go get github.com/jakewright/compile-daemon

EXPOSE 80

WORKDIR /app
COPY . .

RUN go mod download

# In order for a build argument to be available in the CMD, we must make it an
# environment variable. This is because the CMD is only executed at runtime. The
# ARG command must be after FROM to be available at this point in the Dockerfile.
ARG service_name
ENV SERVICE ${service_name}

# Copy assets for this service into /assets in the image. The LICENCE file is
# included because Docker requires at least one file in the COPY command, and
# it is assumed that LICENCE will always exist. Any file could be used.
COPY LICENCE ./private/assets/${service_name}/dev/* /assets/

# Must use exec form so that compile-daemon receives signals. The graceful-kill
# option then forwards them to the go binary. The -directories option doesn't
# work with the directories the other way around. It might be because of the dot
# in the service name.
CMD ["sh", "-c", "compile-daemon -build=\"go install ./services/${SERVICE}\" -command=/go/bin/${SERVICE} -directories=libraries/go,services/${SERVICE} -log-prefix=false -log-prefix=false -graceful-kill=true -graceful-timeout=10"]
