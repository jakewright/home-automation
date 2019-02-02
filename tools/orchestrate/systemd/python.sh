#!/bin/bash

# Abort if anything fails
set -e

if [ ! -d "env" ]; then
    echo "Creating virtualenv"
    virtualenv env
fi

echo "Activating virtualenv"
source env/bin/activate
echo "Installing requirements"
env/bin/pip install -r requirements.txt
deactivate
