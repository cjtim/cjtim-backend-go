#!/bin/sh

go test -coverprofile=coverage.txt -covermode count ./...
go tool cover -func=coverage.txt