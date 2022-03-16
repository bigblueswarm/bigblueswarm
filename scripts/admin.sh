#!/bin/bash

API_KEY=$(grep  -Po "apiKey: (.*)" $(dirname "$0")/../config.yml | cut -d: -f2 | xargs)

usage() {
  echo "
    Usage: admin.sh [OPTION]

    Call b3lb admin functions
    -l, --list           Call admin list api. It displays instance list.
    -c, --create         Call admin create api. It prompts to ask the instance url and instance secret to create.
    -d, --delete         Call admin delete api. It prompts to ask the instance url to delete.
    -s, --status         Call admin status api.
    -h, --help           Display usage
  "
}

create() {
    echo "Enter the instance url to create:"
    read INSTANCE_URL
    echo "Enter the instance secret:"
    read INSTANCE_SECRET
    curl -s -X POST http://localhost:8090/admin/servers -H "Authorization: $API_KEY" -H 'Content-Type: application/json' -d "{\"url\":\"$INSTANCE_URL\", \"secret\": \"$INSTANCE_SECRET\"}" | jq "."
}

delete() {
    echo "Enter the instance url to delete:"
    read INSTANCE_URL
    curl -G -X DELETE http://localhost:8090/admin/servers -H "Authorization: $API_KEY" --data-urlencode "url=$INSTANCE_URL" | jq "."
}

list() {
    curl -s -G -X GET http://localhost:8090/admin/servers -H "Authorization: $API_KEY" | jq "."
}

cluster_status() {
  curl -s -G -X GET http://localhost:8090/admin/cluster/status -H "Authorization: $API_KEY" | jq "."
}

for param in "$@"
do
  case $param in
    -c | --create)
      create
      ;;
    -l | --list)
      list
      ;;
    -d | --delete)
      delete
      ;;
    -s | --status)
      cluster_status
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
