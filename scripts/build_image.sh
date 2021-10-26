#!/bin/bash
DOCKER_BUILDKIT=0 docker build "$(dirname "$0")/docker" -t sledunois/bbb-dev:2.4-develop