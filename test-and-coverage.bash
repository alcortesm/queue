#!/usr/bin/env bash

# This shell script runs the tests (with race detection) in all
# packages, generates a report for each of them individually and
# consolidates all the reports in a single coverage.txt file in the root
# of the repo.

set -e

REPORT="coverage.txt"

echo "" > ${REPORT}
for pkg in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $pkg
    if [ -f profile.out ]; then
        cat profile.out >> ${REPORT}
        rm profile.out
    fi
done
