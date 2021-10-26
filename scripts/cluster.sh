#!/bin/bash

start() {
    echo "Starting BigBlueButton cluster..."
    docker-compose -f "$(dirname "$0")/docker-compose.yml" up -d
}

stop() {
    echo "Stopping cluster..."
    docker rm -f bbb1 bbb2 influxdb
}

usage() {
  echo "
    Usage: cluster.sh [OPTION]

    Manage local development cluster.
    -r, --run, --start        Start development cluster. It starts an InfluxDB serveur and two BigBlueButton servers.
    -s, --stop                Stop development cluster.
    -t, --set-token [token]   Update BigBlueButton server to use given InfluxDB token.
    -h, --help                Display this help and exit
  "
}

set_influxdb_token() {
  echo "Setting up InfluxDB token..."
  docker exec -it bbb1 sh -c "echo INFLUXDB_TOKEN=$1 > /etc/default/telegraf && . /etc/default/telegraf && systemctl restart telegraf"
  docker exec -it bbb2 sh -c "echo INFLUXDB_TOKEN=$1 > /etc/default/telegraf && . /etc/default/telegraf && systemctl restart telegraf"
  echo "Done"
}

for param in "$@"
do
  case $param in
    -r | --run | --start)
      start
      ;;
    -s | --stop)
      stop
      ;;
    -t | --set-token)
      set_influxdb_token "$2"
      exit 1
      ;;
    -h | --help)
      usage
      exit 1
      ;;
    *)
      echo "Invalid argument : $param"
  esac
  if [ ! $? -eq 0 ]; then
    exit 1
  fi
done