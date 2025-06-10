#!/usr/bin/env bash

# Build the frontend assets and start the Go backend.
# Passing -n skips the frontend build step.

set -e

# parse options
BUILD_FRONTEND=true
usage() {
  echo "Usage: $0 [-n]" >&2
  echo "  -n  do not build frontend" >&2
}

while getopts "hn" opt; do
  case $opt in
    h)
      usage
      exit 0
      ;;
    n)
      BUILD_FRONTEND=false
      ;;
  esac
done
shift $((OPTIND-1))

if $BUILD_FRONTEND; then
  echo "Building frontend..."
  (cd src/frontend && npm i && npm run build)
fi

echo "Starting backend..."
go run ./src/backend