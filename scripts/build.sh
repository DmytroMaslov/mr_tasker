#!/bin/bash
appname=$1
go build -o ./tmp/bin/$appname ./cmd/$appname/main.go