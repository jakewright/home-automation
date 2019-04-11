#!/bin/bash

# Abort if anything fails
set -e

DASHES=$(echo $SERVICE | tr "." "-")

ssh -t -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null "$TARGET_USERNAME"@"$DEPLOYMENT_TARGET" << EOF
    # Abort if anything fails
    set -e

    cd $TARGET_DIRECTORY/src

    REV=\$(git log --pretty=format:'%h' -n 1)
    echo "Current revision: \$REV"

    sudo systemctl status $DASHES.service
EOF
