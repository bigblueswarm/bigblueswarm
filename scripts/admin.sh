#!/bin/bash

API_KEY=$(grep  -Po "api_key: (.*)" $(dirname "$0")/../config.yml | cut -d: -f2 | xargs)

usage() {
  echo "
    Usage: admin.sh [OPTION]

    Call b3lb admin functions
    -l, --list           Call admin list api. It displays instance list.
    -c, --create         Call admin create api. It prompts to ask the instance url and instance secret to create.
    -d, --delete         Call admin delete api. It prompts to ask the instance url to delete.
    -h, --help           Display usage
  "
}

create() {
    echo "Enter the instance url to create:"
    read INSTANCE_URL
    echo "Enter the instance secret:"
    read INSTANCE_SECRET
    http POST :8090/admin/servers url="$INSTANCE_URL" secret="$INSTANCE_SECRET" "Authorization: $API_KEY"
}

delete() {
    echo "Enter the instance url to delete:"
    read INSTANCE_URL
    http DELETE :8090/admin/servers url=="$INSTANCE_URL" "Authorization: $API_KEY"
}

list() {
  http GET :8090/admin/servers "Authorization: $API_KEY"
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
