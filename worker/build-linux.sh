#!/bin/bash

GOOS=linux GOARCH=amd64 go build cmd/worker/main.go
