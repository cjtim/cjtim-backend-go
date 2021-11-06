#!/bin/sh

go test -v -coverprofile=coverage.out ./...

# tail -q -n +2 coverage.out >> cover/coverage.cov
go tool cover -func=coverage.out