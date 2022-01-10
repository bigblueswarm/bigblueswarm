#!/bin/bash

echo "$(grep -Po "version = \"(.*)\"" $(dirname "$0")/../cmd/b3lb/main.go | cut -d= -f2 | xargs)"