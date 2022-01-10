#!/bin/bash
GREEN="\e[32m"
MAGENTA="\e[35m"
ENDCOLOR="\e[0m"

VERSION="$(sh -c "./$(dirname "$0")/version.sh")"

set -e

dash() {
    echo "$(/usr/bin/printf "\xE2\x9C\x94")"
}

log() {
    echo -e "${MAGENTA}[BUILD]${ENDCOLOR} $1"
}

log "Building B3LB artifact"
rm -rf "$(dirname "$0")"/../bin/
go build -o "$(dirname "$0")"/../bin/"b3lb-$VERSION" "$(dirname "$0")"/../cmd/b3lb/main.go
log "Done ${GREEN}$(dash)${ENDCOLOR}"