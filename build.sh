#!/usr/bin/env bash

go test ./...

env GOOS=darwin GOARCH=amd64 go build -o feedback-service-mac
env GOOS=linux GOARCH=amd64 go build -o feedback-service-linux
env GOOS=windows GOARCH=amd64 go build