#!/bin/bash

# Abort if anything fails
set -e

if [ ! -f "$1"/deploy.cfg ]; then
    echo "Configuration file not found for $1."
    exit 1
fi

source "$1"/deploy.cfg

echo "Deploying $SERVICE_NAME to $DEPLOYMENT_TARGET..."

ssh -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null "$TARGET_USERNAME"@"$DEPLOYMENT_TARGET" << EOF

    cd $PROJECT_PARENT

    if [ ! -d "home-automation" ]; then
        git clone https://github.com/jakewright/home-automation.git
    fi

    cd home-automation

    git checkout -- .
    git checkout master
    git pull

    if [ "$INIT_SYSTEM" = "systemd" ]; then
        sudo systemctl stop "$SERVICE_NAME".service
    fi

    cd "$1"

    bash ./deploy.sh

    if [ "$INIT_SYSTEM" = "systemd" ]; then
        echo "Creating systemd service"
        # The quotes are needed around the variable to preserve the new lines
        echo "$SYSTEMD_SERVICE" | sudo tee /lib/systemd/system/$SERVICE_NAME.service
        sudo chmod 644 /lib/systemd/system/$SERVICE_NAME.service

        sudo systemctl daemon-reload
        sudo systemctl enable $SERVICE_NAME.service
        sudo systemctl start $SERVICE_NAME.service
    fi
EOF
