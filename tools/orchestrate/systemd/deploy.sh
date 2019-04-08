#!/bin/bash

# Abort if anything fails
set -e

DASHES=$(echo $SERVICE | tr "." "-")

ssh -t -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null "$TARGET_USERNAME"@"$DEPLOYMENT_TARGET" << EOF

    cd $TARGET_DIRECTORY

    if [ ! -d "src" ]; then
        git clone https://github.com/jakewright/home-automation.git src
    fi

    cd src

    git checkout -- .
    git checkout master
    git pull

    sudo systemctl stop "$DASHES".service

    cd $SERVICE
    bash ../tools/orchestrate/systemd/$LANG.sh

    echo "Creating systemd service"
    # The quotes are needed around the variable to preserve the new lines
    echo "$SYSTEMD_SERVICE" | sudo tee /lib/systemd/system/$DASHES.service
    sudo chmod 644 /lib/systemd/system/$DASHES.service

    sudo systemctl daemon-reload
    sudo systemctl enable "$DASHES".service
    sudo systemctl start "$DASHES".service

EOF
