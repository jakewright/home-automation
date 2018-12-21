#!/bin/bash

# Abort if anything fails
set -e

printf "Formatting code..."
go fmt $(go list ./... | grep -v /vendor/)

printf "\nRunning golint..."
golint $(go list ./... | grep -v /vendor/)

printf "\nRunning go vet..."
go vet $(go list ./... | grep -v /vendor/)