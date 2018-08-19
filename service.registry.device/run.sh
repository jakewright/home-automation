#!/bin/bash

export FLASK_APP=device_registry
export FLASK_ENV=development
export FLASK_RUN_HOST=0.0.0.0
export FLASK_RUN_PORT=80
flask run
