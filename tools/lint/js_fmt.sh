#!/bin/bash

printf "Running Prettier..."
prettier --write **/*.js

printf "\nRunning eslint..."
eslint --ext .vue --fix ./