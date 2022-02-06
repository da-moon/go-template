#!/usr/bin/env bash
set -e -o pipefail
exec 5>&1
output="$(gofmt -l -w "$@" | tee /dev/fd/5)"
[[ -z "$output" ]]
