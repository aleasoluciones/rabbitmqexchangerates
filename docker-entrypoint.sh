#!/bin/bash
set -e

# Inspired on https://github.com/docker-library/postgres script

case "$1" in
    "rabbitmqexchangerates")
        shift
        OPTIONS="$@"
        exec /go/bin/app ${OPTIONS}
    ;;
esac
exec "$@"
