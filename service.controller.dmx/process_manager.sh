#!/bin/bash

# Start OLA
nohup sh -c /start.sh &

export FLASK_APP=dmx
export FLASK_ENV=development
export FLASK_RUN_HOST=0.0.0.0
export FLASK_RUN_PORT=5006

# Wait for OLA to respond and then start the DMX service
/wait-for-it.sh localhost:9090 -- flask run
