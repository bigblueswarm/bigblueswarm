#!/bin/bash
BLUE="\e[34m"
GREEN="\e[32m"
ENDCOLOR="\e[0m"
CLUSTER_SCRIPT="$(dirname "$0")"/cluster.sh
VERSION="$(sh -c "./$(dirname "$0")/version.sh")"

set -e

dash() {
    echo "$(/usr/bin/printf "\xE2\x9C\x94")"
}

log() {
    echo -e "${BLUE}[INTEGRATION TESTS]${ENDCOLOR} $1"
}

end () {
    echo -e "${GREEN}Done${ENDCOLOR}"
}

# Building artifact
build() {
    sh -c "$(dirname "$0")/build.sh";
}

# Launching integration tests using newman
launch(){
    log "Getting BigBlueButton secret"
    SECRET="$(docker exec bbb1 sh -c "bbb-conf --secret" | grep -Po "Secret: (.*)" | cut -d: -f2 | xargs)"
    log "Launching integration tests"
    npm install newman
    "$(dirname "$0")"/../node_modules/.bin/newman run "$(dirname "$0")"/../test/B3LB.postman_collection.json -e "$(dirname "$0")"/../test/Integration\ test.postman_environment.json --env-var instance_secret="$SECRET" --bail --verbose --ignore-redirects
}

# Starting cluster
start() {
    # Check if the bigbluebutton image already exists. If not, build it.
    if [ "$(docker images | grep sledunois/bbb-dev | wc -l)" -eq "0" ]; then
        log "BigBlueButton image not found, building it..."
        sh -c "$(dirname "$0")/build_image.sh";
        log "Image built ${GREEN}$(dash)${ENDCOLOR}"
    fi

    # Launching cluster and setting up the cluster configuration
    log "Launching cluster"
    sh -c "./$CLUSTER_SCRIPT -r"
    sleep 5m
    log "Setting up cluster configuration"
    sh -c "./$CLUSTER_SCRIPT -i"
    log "Starting B3LB artifact $VERSION"
    nohup ./$(dirname "$0")/../bin/b3lb-$VERSION --config config.yml &
    echo $! > "$(dirname "$0")"/../bin/b3lb.pid
    sleep 15s
    end
}

# Stopping cluster
stop() {
    log "Killing B3LB process"
    kill -9 $(cat "$(dirname "$0")"/../bin/b3lb.pid)
    rm -f "$(dirname "$0")"/../bin/b3lb.pid
    end
    log "Stopping cluster"
    sh -c "./$CLUSTER_SCRIPT -s"
    log "Removing cluster containers"
    docker rm -f bbb1 bbb2 redis influxdb
    end
}

usage() {
  echo "
    Usage: integration_test.sh [OPTION]

    Manage local development cluster.
    -r, --run, --start          Start integration test cluster. It starts an InfluxDB serveur and two BigBlueButton servers.
    -s, --stop                  Stop integration test cluster.
    -b, --build                 Build the BigBlueButton artifact.
    -l, --launch                Launch the integration test.
    -h, --help                  Display this help and exit
  "
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
    -b | --build)
      build
      ;;
    -l | --launch)
      launch
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