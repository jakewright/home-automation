#!/bin/bash

# Start OLA
nohup sh -c /start.sh &

# Wait for OLA to respond and then start the DMX service
/wait-for-it.sh localhost:9090 -- python /usr/src/app/run.py
