#!/bin/bash

start() {
    docker-compose -f "$(dirname "$0")/docker-compose.yml" up -d
}

stop() {
    docker rm -f bbb1 bbb2 influxdb
}

usage() {
  echo "Help!"
}

for param in "$@"
do
  case $param in
    -run | --start)
      start
      ;;
    --stop)
      stop
      ;;
    -h | --help)
      usage
      ;;
    *)
      echo "Invalid argument : $param"
  esac
  if [ ! $? -eq 0 ]; then
    exit 1
  fi
done