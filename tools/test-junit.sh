#!/bin/sh

go get gotest.tools/gotestsum
gotestsum --junitfile report.xml --format testname