#!/bin/sh
# export MONGO_URI=mongodb://mongodb:mongodb@localhost:27017
go test -v -coverprofile=coverage.out ./...

# tail -q -n +2 coverage.out >> cover/coverage.cov
go tool cover -func=coverage.out