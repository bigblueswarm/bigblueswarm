#!/bin/bash

TOKEN=$(grep  -Po "token: (.*)" $(dirname "$0")/../config.yml | cut -d: -f2  | xargs)

start() {
    echo "Starting BigBlueButton cluster..."
    docker-compose -f "$(dirname "$0")/docker-compose.yml" up -d
}

stop() {
    echo "Stopping cluster..."
    docker stop bbb1 bbb2 influxdb redis
}

usage() {
  echo "
    Usage: cluster.sh [OPTION]

    Manage local development cluster.
    -r, --run, --start          Start development cluster. It starts an InfluxDB serveur and two BigBlueButton servers.
    -s, --stop                  Stop development cluster.
    -i, --init, --init-cluster  Initialize the cluster.
    -t, --set-token [token]     Update BigBlueButton server to use given InfluxDB token.
    -h, --help                  Display this help and exit
  "
}

set_influxdb_token() {
  echo "Setting up InfluxDB token..."
  docker exec -it bbb1 sh -c "echo 'INFLUXDB_TOKEN=$1\nB3LB_HOST=http://localhost/bigbluebutton' > /etc/default/telegraf && . /etc/default/telegraf && systemctl restart telegraf"
  docker exec -it bbb2 sh -c "echo 'INFLUXDB_TOKEN=$1\nB3LB_HOST=http://localhost:8080/bigbluebutton' > /etc/default/telegraf && . /etc/default/telegraf && systemctl restart telegraf"
  echo "Done"
}

init_cluster() {
  echo "Initializing cluster..."
  docker exec -it influxdb sh -c "influx setup --name b3lbconfig --org b3lb --username admin --password password --token ${TOKEN} --bucket bucket --retention 0 --force"
  set_influxdb_token "$TOKEN"
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
    -i | --init | --init-cluster)
      init_cluster
      ;;
    -t | --set-token)
      set_influxdb_token "$TOKEN"
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
