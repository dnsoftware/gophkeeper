#!/bin/bash

GOOS=windows GOARCH=amd64 /usr/local/go/bin/go build -o clientkeeper_win.exe -ldflags="-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date)' -X main.buildCommit=$(git rev-parse HEAD)"