#!/bin/bash

./serverkeeper -s=/home/dmitry/go/src/gophkeeper/cmd/server/certs/server.crt -p=/home/dmitry/go/src/gophkeeper/cmd/server/certs/server.key -c=/home/dmitry/go/src/gophkeeper/cmd/server/config.yaml -d=postgres://p1pool:pass@localhost:54321/gophkeeper?sslmode=disable