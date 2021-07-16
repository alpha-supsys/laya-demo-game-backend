#!/bin/sh

target=main

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/$target ./src/main.go