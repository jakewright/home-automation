#!/bin/bash

# Abort if anything fails
set -e

UNDERSCORES=$(echo $SERVICE | tr "." "_" | tr "-" "_")
DASHES=$(echo $SERVICE | tr "." "-")

ssh -t -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null "$TARGET_USERNAME"@"$DEPLOYMENT_TARGET" << EOF
    # Abort if anything fails
    set -e

    echo "Building $SERVICE..."
    cd $TARGET_DIRECTORY/src
    git pull

    # Get the current git commit hash to use as the image label
    HASH=\$(git log --pretty=format:'%h' -n 1)

    # Use a file called prod.dockerfile if it exists, otherwise use Dockerfile
    DOCKER_FILE=prod.dockerfile
    if [ ! -f "./$SERVICE/\$DOCKER_FILE" ]; then
        DOCKER_FILE=Dockerfile
    fi

    # Build and push the Docker image
    # Escape the variables that are defined within the SSH context, otherwise
    # the shell will try to replace them with variables defined locally.
    docker build -f ./$SERVICE/\$DOCKER_FILE -t localhost:6000/jakewright/home-automation-$DASHES:\$HASH --rm .
    docker push localhost:6000/jakewright/home-automation-$DASHES:\$HASH

    echo "Pulling new image..."
    cd /volume1/docker/home-automation/
    docker pull localhost:6000/jakewright/home-automation-$DASHES:\$HASH

    cd $TARGET_DIRECTORY
    echo $UNDERSCORES"_VERSION=\$HASH" >> .env

    docker-compose stop $SERVICE
    docker-compose up -d $SERVICE
EOF