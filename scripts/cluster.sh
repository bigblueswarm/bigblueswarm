#!/bin/bash
set -e 
YELLOW="\e[33m"
GREEN="\e[32m"
ENDCOLOR="\e[0m"

TOKEN=$(grep  -Po "token: (.*)" $(dirname "$0")/../config.yml | cut -d: -f2  | xargs)

log() {
    echo -e "${YELLOW}[CLUSTER]${ENDCOLOR} $1"
}

dash() {
    echo "$(/usr/bin/printf "\xE2\x9C\x94")"
}

start() {
    log "Starting BigBlueButton cluster"
    docker-compose -f "$(dirname "$0")/docker-compose.yml" up -d
    log "Started ${GREEN}$(dash)${ENDCOLOR}"
}

stop() {
    log "Stopping cluster"
    docker stop bbb1 bbb2 influxdb redis
    log "Stopped ${GREEN}$(dash)${ENDCOLOR}"
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
  log "Setting up InfluxDB token"
  docker exec bbb1 sh -c "echo 'INFLUXDB_TOKEN=$1\nB3LB_HOST=http://localhost/bigbluebutton' > /etc/default/telegraf && . /etc/default/telegraf && systemctl restart telegraf"
  docker exec bbb2 sh -c "echo 'INFLUXDB_TOKEN=$1\nB3LB_HOST=http://localhost:8080/bigbluebutton' > /etc/default/telegraf && . /etc/default/telegraf && systemctl restart telegraf"
  log "Done ${GREEN}$(dash)${ENDCOLOR}"
}

init_cluster() {
  log "Initializing cluster"
  log "Setting up InfluxDB token"
  docker exec influxdb sh -c "influx setup --name b3lbconfig --org b3lb --username admin --password password --token ${TOKEN} --bucket bucket --retention 0 --force"
  log "${GREEN}Done${ENDCOLOR}"
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
      log "Invalid argument : $param"
  esac
  if [ ! $? -eq 0 ]; then
    exit 1
  fi
done
