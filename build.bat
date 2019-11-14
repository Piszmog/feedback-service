@echo off

set GOOS=windows
set GOARCH=amd64
go test ./...
go build

set GOOS=linux
go build -o feedback-service-linux

set GOOS=darwin
go build -o feedback-service-mac