#!/bin/bash

# Abort if anything fails
set -e

UNDERSCORES=$(echo $SERVICE | tr "." "_" | tr "-" "_")

ssh -t -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null "$TARGET_USERNAME"@"$DEPLOYMENT_TARGET" << EOF
    # Abort if anything fails
    set -e

    cd $TARGET_DIRECTORY
    REV=\$(cat .env | grep $UNDERSCORES | tail -1 | cut -d "=" -f 2)
    echo "Current revision: \$REV"
    docker ps --filter "name=$SERVICE" --format "{{.Image}} {{.Status}}"
EOF