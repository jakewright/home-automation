#!/bin/bash

# Abort if anything fails
set -e

UNDERSCORES=$(echo $SERVICE | tr "." "_" | tr "-" "_")
DASHES=$(echo $SERVICE | tr "." "-")

ssh -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null "$TARGET_USERNAME"@"$DEPLOYMENT_TARGET" << EOF
    echo "Building $SERVICE..."
    cd $TARGET_DIRECTORY/src
    git pull
    HASH=\$(git log --pretty=format:'%h' -n 1)
    docker build -f ./$SERVICE/prod.dockerfile -t localhost:6000/jakewright/home-automation-$DASHES:\$HASH --rm .
    docker push localhost:6000/jakewright/home-automation-$DASHES:\$HASH

    echo "Pulling new image..."
    cd /volume1/docker/home-automation/
    docker pull localhost:6000/jakewright/home-automation-$DASHES:\$HASH

    cd $TARGET_DIRECTORY
    echo $UNDERSCORES"_VERSION=\$HASH" >> .env

    docker-compose stop $SERVICE
    docker-compose up -d $SERVICE
EOF